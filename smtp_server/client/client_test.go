package client

import (
  "errors"
  "testing"
  "github.com/stretchr/testify/assert"
)

type legitRWC struct {
  bytesWritten int
}

func (l *legitRWC) Read(b []byte) (int, error) {
  copy(b, []byte("foobar"))
  return 6, nil
}

func (l *legitRWC) Write(b []byte) (int, error) {
  l.bytesWritten += len(b)
  return len(b), nil
}

func (l *legitRWC) Close() error {
  return nil
}

type errRWC struct {
  bytesWritten int
}

func (l *errRWC) Read(_ []byte) (int, error) {
  return 0, errors.New("whoa a read err!")
}

func (l *errRWC) Write(_ []byte) (int, error) {
  return 0, errors.New("whoa a write err!")
}

func (l *errRWC) Close() error {
  return errors.New("whoa a close err!")
}

func TestGoodReadPassesThrough(t *testing.T) {
  c := MakeClient(new(legitRWC))
  got := make([]byte, 6)
  n, err := c.in.Read(got)
  assert.Equal(t, 6, n)
  assert.NoError(t, err)
  assert.Equal(t, []byte("foobar"), got)
}

func TestGootWritePassesThrough(t *testing.T) {
  c := MakeClient(new(legitRWC))
  n, err := c.out.Write([]byte("foobar"))
  assert.NoError(t, err)
  assert.Equal(t, 6, n)
}

func TestBadReadPassesThrough(t *testing.T) {
  c := MakeClient(new(errRWC))
  got := make([]byte, 6)
  n, err := c.in.Read(got)
  assert.Equal(t, 0, n)
  assert.Error(t, err)
  assert.Equal(t, "whoa a read err!", err.Error())
  assert.Equal(t, make([]byte, 6), got)
}

func TestBadWritePassesThrough(t *testing.T) {
  c := MakeClient(new(errRWC))
  c.out.Write([]byte("foobar"))
  err := c.out.Flush()
  assert.Error(t, err)
  assert.Equal(t, "whoa a write err!", err.Error())
}
