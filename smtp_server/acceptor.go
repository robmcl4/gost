package smtp_server

import (
  "net"
  "fmt"
  log "github.com/Sirupsen/logrus"
  "github.com/robmcl4/gost/config"
)

func getServerConnection() (net.Listener, error) {
  addr := fmt.Sprintf("%s:%d",
                      config.GetListenAddress(),
                      config.GetListenPort())
  return net.Listen("tcp", addr)
}

func listenForConnections(l net.Listener, c chan net.Conn) error {
  listeningAddr := "<nil>"
  if l.Addr() != nil {
    listeningAddr = l.Addr().String()
  }
  log.WithFields(log.Fields{
    "listening_on": listeningAddr,
    "fqdn": config.GetFQDN(),
  }).Info("Starting connection listener")
  for {
    conn, err := l.Accept()
    if err != nil {
      return err
    }
    c <- conn
  }
}
