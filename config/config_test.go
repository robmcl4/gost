package config

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestDefaultPort(t *testing.T) {
  assert.Equal(t, 587, GetListenPort(), "default port should be 587")
}


func TestDefaultHost(t *testing.T) {
  assert.Equal(
    t,
    "127.0.0.1",
    GetListenAddress(),
    "default address should be 127.0.0.1")
}
