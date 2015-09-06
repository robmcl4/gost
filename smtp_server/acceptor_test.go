package server

import (
  "testing"
  "net"
  "errors"
  "github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------

func TestGetServerConnection(t *testing.T) {
  l, err := getServerConnection()
  if err != nil {
    t.Errorf("Error making server connection")
    return
  }
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
  c := make(chan bool)
  h := func(conn net.Conn) {
    c <- true
  }
  err := listenForConnections(&mylistener{}, h)
  assert.Equal(t, "myerror", err.Error(), "error message should be \"foo\"")
  <- c
}
