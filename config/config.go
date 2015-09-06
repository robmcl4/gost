// Contains configuration loading behavior and defaults for gost.
package config

import (
  "sync"
)

type configuration struct {
  listenAddress string
  listenPort    int
  fqdn          string
}

// The global configuration state, set with defaults.
var globalConfig configuration = configuration{ "127.0.0.1",
                                                587,
                                                "mail.example.com" }

// The read-write lock for ensuring safe access to globalConfig
var globalConfigLock sync.RWMutex = sync.RWMutex{}


// Gets the address the server should listen on, for example "127.0.0.1".
func GetListenAddress() string {
  globalConfigLock.RLock()
  ret := globalConfig.listenAddress
  globalConfigLock.RUnlock()
  return ret
}


// Gets the port the server should listen on
func GetListenPort() int {
  globalConfigLock.RLock()
  ret := globalConfig.listenPort
  globalConfigLock.RUnlock()
  return ret
}


func GetFQDN() string {
  globalConfigLock.RLock()
  ret := globalConfig.fqdn
  globalConfigLock.RUnlock()
  return ret
}
