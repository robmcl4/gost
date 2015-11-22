package config

import (
  "time"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestDefaultListenParams(t *testing.T) {
  addr, port := GetListenParams()
  assert.Equal(t, "127.0.0.1", addr)
  assert.Equal(t, 587, port)
}

func TestDefaultFQDN(t *testing.T) {
  assert.Equal(
    t,
    "mail.example.com",
    GetFQDN(),
    "default FQDN should be mail.example.com")
}

func TestDefaultEmailTTL(t *testing.T) {
  assert.Equal(
    t,
    15*60*time.Second,
    GetEmailTTL(),
    "default email TTL should be 15min")
}

func TestDefaultMatcherTTL(t *testing.T) {
  assert.Equal(
    t,
    15*60*time.Second,
    GetMatcherTTL(),
    "default email TTL should be 15min")
}

func TestDefaultBackendType(t *testing.T) {
  assert.Equal(
    t,
    "memory",
    GetBackendType(),
  )
}

func TestSetListenParams(t *testing.T) {
  old_intf, old_port := GetListenParams()
  defer func() {
    globalConfig.listenPort = old_port
    globalConfig.listenAddress = old_intf
  }()
  SetListenParams("192.168.1.1", 11211)
  addr, port := GetListenParams()
  assert.Equal(t, 11211, port)
  assert.Equal(t, "192.168.1.1", addr)
}
