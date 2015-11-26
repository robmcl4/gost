package config

import (
  "strings"
  "os"
  "io/ioutil"
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

func TestMissingFileGivesError(t *testing.T) {
  teardownConfigFile()
  _, err := detectFileLocation()
  assert.Error(t, err)
}

func TestMissingFileCannotLoad(t *testing.T) {
  teardownConfigFile()
  oldConfig := globalConfig
  err := LoadConfigFromFile()
  assert.Error(t, err)
  assert.Equal(t, oldConfig, globalConfig)
}

func TestUnreadableFileCannotLoad(t *testing.T) {
  teardownConfigFile()
  assert.NoError(t, ioutil.WriteFile("config.json", []byte("{}"), 0222))
  oldConfig := globalConfig
  err := LoadConfigFromFile()
  assert.Error(t, err)
  assert.Equal(t, oldConfig, globalConfig)
  assert.NoError(t, os.Remove("config.json"))
}

func TestNonJsonCannotLoad(t *testing.T) {
  teardownConfigFile()
  assert.NoError(t, ioutil.WriteFile("config.json", []byte("uhhh i'm not json"), 0644))
  oldConfig := globalConfig
  err := LoadConfigFromFile()
  assert.Error(t, err)
  assert.Equal(t, oldConfig, globalConfig)
  assert.NoError(t, os.Remove("config.json"))
}

func TestFindsFile(t *testing.T) {
  assert.NoError(t, putConfigFile())
  defer teardownConfigFile()
  loc, err := detectFileLocation()
  assert.True(t, strings.HasSuffix(loc, "config.json"), loc)
  assert.NoError(t, err)
}

func TestLoadsBasicFile(t *testing.T) {
  assert.NoError(t, putConfigFile())
  defer teardownConfigFile()
  oldConfig := globalConfig
  defer func() {
    globalConfig = oldConfig
  }()

  err := LoadConfigFromFile()
  assert.NoError(t, err)
}

func putConfigFile() error {
  return ioutil.WriteFile(
    "config.json",
    []byte(
`{
  "listen_address": "1.1.1.1",
  "listen_port": 12345,
  "fqdn": "foobar.example.com",
  "email_ttl": 11,
  "matcher_ttl": 11,
  "backendType": "memory"
}
`,
    ),
    0644,
  )
}

func teardownConfigFile() error {
  return os.Remove("config.json")
}
