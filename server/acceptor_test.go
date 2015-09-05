package server

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

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
