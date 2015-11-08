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
  backend       string
}

// The global configuration state, set with defaults.
var globalConfig configuration = configuration{
  listenAddress: "127.0.0.1",
  listenPort:    587,
  fqdn:          "mail.example.com",
  emailTtl:      15*60*time.Second,
  matcherTtl:    15*60*time.Second,
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
  return globalConfig.backend
}

func SetListenParams(intrfce string, port int) {
  globalConfig.Lock()
  defer globalConfig.Unlock()
  globalConfig.listenAddress = intrfce
  globalConfig.listenPort = port
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
  conf := make(map[string]interface{})
  if err = json.Unmarshal(bindata, &conf); err != nil {
    return err
  }

  newConfig := configuration{}
  if val, err := get_string(conf, "listen_address"); err == nil {
    newConfig.listenAddress = val
  } else {
    return err
  }
  if val, err := get_int(conf, "listen_port"); err == nil {
    newConfig.listenPort = val
  } else {
    return err
  }
  if val, err := get_string(conf, "fqdn"); err == nil {
    newConfig.fqdn = val
  } else {
    return err
  }
  if val, err := get_int(conf, "email_ttl"); err == nil {
    newConfig.emailTtl = time.Duration(val)*time.Minute
  } else {
    return err
  }
  if val, err := get_int(conf, "matcher_ttl"); err == nil {
    newConfig.matcherTtl = time.Duration(val)*time.Minute
  } else {
    return err
  }
  if val, ok := conf["backend"]; ok {
    if m, ok := val.(map[string]interface{}); ok {
      if backend, err := get_string(m, "type"); err == nil {
        newConfig.backend = backend
      } else {
        return err
      }
    } else {
      return errors.New("backend is not a string")
    }
  } else {
    return errors.New("no backend section found")
  }

  oldConfig := globalConfig
  oldConfig.Lock()
  globalConfig = newConfig
  oldConfig.Unlock()
  return nil
}

func get_int(m map[string]interface{}, key string) (int, error) {
  if m == nil {
    return 0, errors.New("could find key=\""+key+"\" for nil map")
  }
  if val, ok := m[key]; ok {
    if i, ok := val.(float64); ok {
      return int(i), nil
    } else {
      return 0, errors.New("key was not int key=\""+key+"\"")
    }
  }
  return 0, errors.New("could not find key=\""+key+"\" in map")
}

func get_string(m map[string]interface{}, key string) (string, error) {
  if m == nil {
    return "", errors.New("could find key=\""+key+"\" for nil map")
  }
  if val, ok := m[key]; ok {
    if str, ok := val.(string); ok {
      return str, nil
    } else {
      return "", errors.New("key was not string key=\""+key+"\"")
    }
  }
  return "", errors.New("could not find key=\""+key+"\" in map")
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
