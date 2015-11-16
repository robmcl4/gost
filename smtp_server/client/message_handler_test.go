package client

import (
  "net"
  "bytes"
  "bufio"
  "testing"
  "errors"
  "github.com/robmcl4/gost/email"
  "github.com/stretchr/testify/assert"
)

type myConn struct {
  net.Conn
  closed bool
}

func (m *myConn) Close() error {
  m.closed = true
  return nil
}

type errConn struct {
  net.Conn
}

func (_ *errConn) Close() error {
  return errors.New("My Close Error")
}

type errWriter struct { }

func (_ *errWriter) Write(_ []byte) (int, error) {
  return 0, errors.New("Write error!")
}

type errReader struct { }

func (_ *errReader) Read(_ []byte) (int, error) {
  return 0, errors.New("Read error!")
}

func TestCloseTerminatesConnection(t *testing.T) {
  mybuf := new(bytes.Buffer)
  mycon := myConn{}
  c := Client{
    &mycon,
    bufio.NewReader(new(bytes.Buffer)),
    bufio.NewWriter(mybuf),
  }
  err := c.Close()
  assert.Equal(
    t,
    "421 Service Unavailable: Terminating Connection\r\n",
    mybuf.String(),
  )
  assert.NoError(t, err)
  assert.True(t, mycon.closed, "Close() should have been called")
}

func TestCloseBubblesUpError(t *testing.T) {
  mybuf := new(bytes.Buffer)
  mycon := errConn{}
  c := Client{
    &mycon,
    bufio.NewReader(new(bytes.Buffer)),
    bufio.NewWriter(mybuf),
  }
  err := c.Close()
  assert.Equal(
    t,
    "421 Service Unavailable: Terminating Connection\r\n",
    mybuf.String(),
  )
  assert.Error(t, err)
  assert.Equal(t, err.Error(), "My Close Error")
}

func TestCloseIgnoresWriteErrors(t *testing.T) {
  mycon := myConn{}
  c := Client{
    &mycon,
    bufio.NewReader(new(bytes.Buffer)),
    bufio.NewWriter(new(errWriter)),
  }
  err := c.Close()
  assert.NoError(t, err)
  assert.True(t, mycon.closed, "Close() should have been called")
}

func TestOpensWithServiceReady(t *testing.T) {
  mybuf := new(bytes.Buffer)
  c := Client{
    &myConn{},
    bufio.NewReader(new(bytes.Buffer)),
    bufio.NewWriter(mybuf),
  }
  ch := make(chan *email.SMTPEmail, 1)
  assert.Error(t, c.BeginReceive(ch))
  assert.Len(t, ch, 0)
  assert.Equal(t, "220 mail.example.com ESMTP\r\n", mybuf.String())
}

func TestHelo(t *testing.T) {
  output := new(bytes.Buffer)
  input := bytes.NewBufferString("HELO\r\n")
  c := Client{
    &myConn{},
    bufio.NewReader(input),
    bufio.NewWriter(output),
  }
  ch := make(chan *email.SMTPEmail, 1)
  assert.Error(t, c.BeginReceive(ch))
  assert.Len(t, ch, 0)
  assert.Equal(
    t,
    "220 mail.example.com ESMTP\r\n" +
    "250 Ok\r\n",
    output.String(),
  )
}
