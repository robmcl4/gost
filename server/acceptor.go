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

func listenForConnections(l net.Listener, h func(net.Conn)) error {
  for {
    conn, err := l.Accept()
    if err != nil {
      return err
    }
    go h(conn)
  }
}
