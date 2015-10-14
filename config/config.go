// Contains configuration loading behavior and defaults for gost.
package config

import (
  "time"
  "sync"
)

type configuration struct {
  sync.RWMutex
  listenAddress string
  listenPort    int
  fqdn          string
  email_ttl     time.Duration
  matcher_ttl   time.Duration
  backend       string
}

// The global configuration state, set with defaults.
var globalConfig configuration = configuration{
  listenAddress: "127.0.0.1",
  listenPort:    587,
  fqdn:          "mail.example.com",
  email_ttl:     15*60*time.Second,
  matcher_ttl:   15*60*time.Second,
  backend:       "memory",
}


// Gets the address the server should listen on, for example "127.0.0.1".
func GetListenAddress() string {
  globalConfig.RLock()
  defer globalConfig.RUnlock()
  return globalConfig.listenAddress
}


// Gets the port the server should listen on
func GetListenPort() int {
  globalConfig.RLock()
  defer globalConfig.RUnlock()
  return globalConfig.listenPort
}


func GetFQDN() string {
  globalConfig.RLock()
  defer globalConfig.RUnlock()
  return globalConfig.fqdn
}


func GetEmailTTL() time.Duration {
  globalConfig.RLock()
  defer globalConfig.RUnlock()
  return globalConfig.email_ttl
}


func GetMatcherTTL() time.Duration {
  globalConfig.RLock()
  defer globalConfig.RUnlock()
  return globalConfig.matcher_ttl
}


func GetBackendType() string {
  globalConfig.RLock()
  defer globalConfig.RUnlock()
  return globalConfig.backend
}
