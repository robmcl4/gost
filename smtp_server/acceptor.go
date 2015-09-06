package smtp_server

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

func listenForConnections(l net.Listener, c chan net.Conn) error {
  fmt.Println("Started listening for connections")
  for {
    conn, err := l.Accept()
    if err != nil {
      return err
    }
    c <- conn
  }
}
