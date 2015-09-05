package server

import (
  "net"
  "fmt"
  "github.com/robmcl4/gost/config"
)

func getServerConnection() (net.Listener, error) {
  addr := fmt.Sprintf("%s:%d",
                      config.GetListenAddress(),
                      config.GetListenPort())
  return net.Listen("tcp", addr)
}
