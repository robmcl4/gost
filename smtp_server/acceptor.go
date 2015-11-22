package smtp_server

import (
  "net"
  "fmt"
  log "github.com/Sirupsen/logrus"
  "github.com/robmcl4/gost/config"
  "github.com/robmcl4/gost/config/shutdown"
)

func getServerConnection() (net.Listener, error) {
  host, port := config.GetListenParams()
  addr := fmt.Sprintf("%s:%d", host, port)
  return net.Listen("tcp", addr)
}

// Listens for connections on the given listener and puts them in the channel.
// Blocks while still receiving connections.
// Returns an error on network problem, or nil if shutdown requested
func listenForConnections(l net.Listener, c chan net.Conn) error {
  id, shutdownRequested := shutdown.AddShutdownListener("Connection listener")
  defer shutdown.RoutineDone(id)

  log.WithFields(log.Fields{
    "listening_on": fmt.Sprintf("%v", l.Addr()),
    "fqdn": config.GetFQDN(),
  }).Info("Starting connection listener")

  // shutdown using the strategy found here http://stackoverflow.com/a/13419724
  quit := false
  go func() {
    <- shutdownRequested
    quit = true
    l.Close()
  }()

  for {
    conn, err := l.Accept()
    if err != nil {
      if quit {
        log.Info("Shutting down connection listener")
        return nil
      }
      l.Close()
      return err
    }
    c <- conn
  }
}
