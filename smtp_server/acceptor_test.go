package smtp_server

import (
  "testing"
  "net"
  "errors"
  "github.com/stretchr/testify/assert"
  "github.com/robmcl4/gost/config"
  "github.com/robmcl4/gost/config/shutdown"
)

// -----------------------------------------------------------------------------

func TestGetServerConnection(t *testing.T) {
  oldAddr, oldPort := config.GetListenParams()
  config.SetListenParams("127.0.0.1", 45432)
  defer config.SetListenParams(oldAddr, oldPort)

  l, err := getServerConnection()
  defer l.Close()
  assert.NotNil(t, l)
  assert.NoError(t, err)
  assert.Equal(t, "tcp", l.Addr().Network())
  assert.Equal(t, "127.0.0.1:45432", l.Addr().String())
}

// -----------------------------------------------------------------------------

type mylistener struct {
  callsToAccept int
  callsUntilError int
  closed bool
}

func (m *mylistener) Accept() (net.Conn, error) {
  m.callsToAccept += 1
  if (m.callsToAccept > m.callsUntilError || m.closed) {
    return nil, errors.New("myerror")
  }
  return &net.TCPConn{}, nil
}

func (m *mylistener) Close() error {
  m.closed = true
  return nil
}

func (m *mylistener) Addr() net.Addr {
  return nil
}

func TestListenForConnection(t *testing.T) {
  c := make(chan net.Conn, 1)
  err := listenForConnections(&mylistener{callsUntilError: 1}, c)
  assert.Equal(t, "myerror", err.Error(), "error message should be \"foo\"")
  <- c
}

func TestListenForConnectionShutsDown(t *testing.T) {
  c := make(chan net.Conn, 1)
  go func() {
    <- c
    shutdown.Shutdown()
  }()
  listenForConnections(&mylistener{callsUntilError: 6}, c)
}
