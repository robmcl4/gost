// Contains configuration loading behavior and defaults for gost.
package config

import (
  "time"
  "sync"
  "errors"
  "os"
  "path/filepath"
  "io/ioutil"
  "encoding/json"
)

type configuration struct {
  sync.RWMutex
  listenAddress string
  listenPort    int
  fqdn          string
  emailTtl      time.Duration
  matcherTtl    time.Duration
  backendType   string
}

// The global configuration state, set with defaults.
var globalConfig configuration = configuration{
  listenAddress: "127.0.0.1",
  listenPort:    587,
  fqdn:          "mail.example.com",
  emailTtl:      15*60*time.Second,
  matcherTtl:    15*60*time.Second,
  backendType:   "memory",
}

// Gets the address the server should listen on, for example "127.0.0.1".
func GetListenParams() (host string, port int) {
  globalConfig.RLock()
  defer globalConfig.RUnlock()
  return globalConfig.listenAddress, globalConfig.listenPort
}

func GetFQDN() string {
  globalConfig.RLock()
  defer globalConfig.RUnlock()
  return globalConfig.fqdn
}

func GetEmailTTL() time.Duration {
  globalConfig.RLock()
  defer globalConfig.RUnlock()
  return globalConfig.emailTtl
}

func GetMatcherTTL() time.Duration {
  globalConfig.RLock()
  defer globalConfig.RUnlock()
  return globalConfig.matcherTtl
}

func GetBackendType() string {
  globalConfig.RLock()
  defer globalConfig.RUnlock()
  return globalConfig.backendType
}

func SetListenParams(intrfce string, port int) {
  globalConfig.Lock()
  defer globalConfig.Unlock()
  globalConfig.listenAddress = intrfce
  globalConfig.listenPort = port
}

// Only used to deserialize JSON config files
type jsonConfig struct {
  listenAddress string
  listenPort    int
  fqdn          string
  emailTtl      float64
  matcherTtl    float64
  backendType   string
}

func (js *jsonConfig) toConfiguration() configuration {
  return configuration{
    listenAddress: js.listenAddress,
    listenPort:    js.listenPort,
    fqdn:          js.fqdn,
    emailTtl:      time.Duration(js.emailTtl)*time.Minute,
    matcherTtl:    time.Duration(js.emailTtl)*time.Minute,
    backendType:   js.backendType,
  }
}

// Loads configuration from an auto-located file.
// This operation is atomic and isolated; if reading from a file fails
// the existing configuration is not touched and an error is thrown.
// If reading is successful, configuration is stored globally and
// nil is returned.
func LoadConfigFromFile() error {
  fname, err := detectFileLocation()
  if err != nil {
    return err
  }
  bindata, err := ioutil.ReadFile(fname)
  if err != nil {
    return err
  }
  conf := new(jsonConfig)
  if err = json.Unmarshal(bindata, conf); err != nil {
    return err
  }

  oldConfig := globalConfig
  oldConfig.Lock()
  defer oldConfig.Unlock()
  globalConfig = conf.toConfiguration()
  return nil
}

// Attempts to detect location of configuration file.
// Returns an absolute location of the file if found. Otherwise, returns
// an error.
func detectFileLocation() (string, error) {
  if _, err := os.Stat("config.json"); os.IsNotExist(err) {
    return "", errors.New("Could not locate configuration file")
  }
  return filepath.Abs("config.json")
}
