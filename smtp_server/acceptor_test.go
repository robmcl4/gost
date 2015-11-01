package smtp_server

import (
  "testing"
  "net"
  "errors"
  "github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------

func TestGetServerConnection(t *testing.T) {
  l, err := getServerConnection()
  assert.NotNil(t, l)
  assert.NoError(t, err)
  assert.Equal(t, "tcp", l.Addr().Network())
  assert.Equal(t, "127.0.0.1:587", l.Addr().String())
  l.Close()
}

// -----------------------------------------------------------------------------

type mylistener struct {
  callsToAccept int
}

func (m *mylistener) Accept() (net.Conn, error) {
  m.callsToAccept += 1
  if (m.callsToAccept >= 2) {
    return nil, errors.New("myerror")
  }
  return &net.TCPConn{}, nil
}

func (m *mylistener) Close() error {
  return nil
}

func (m *mylistener) Addr() net.Addr {
  return nil
}

func TestListenForConnection(t *testing.T) {
  c := make(chan net.Conn, 1)
  err := listenForConnections(&mylistener{}, c)
  assert.Equal(t, "myerror", err.Error(), "error message should be \"foo\"")
  <- c
}
