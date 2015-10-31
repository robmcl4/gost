package smtp_server

import (
  "net"
  "github.com/robmcl4/gost/email"
  "github.com/robmcl4/gost/smtp_server/client"
  "github.com/robmcl4/gost/config/shutdown"
  log "github.com/Sirupsen/logrus"
)

func ReceiveEmail(c chan *email.SMTPEmail) error {
  l, err := getServerConnection()
  if err != nil {
    return err
  }
  cxnChan := make(chan net.Conn, 10)
  go handleClients(cxnChan, c)
  return listenForConnections(l, cxnChan)
}

func handleClients(cxnChan chan net.Conn, emChan chan *email.SMTPEmail) {
  id, shutdownRequested := shutdown.AddShutdownListener("Client handler")
  defer shutdown.RoutineDone(id)

  for {
    select {
    case cxn := <- cxnChan:
      go handleClient(cxn, emChan)
    case <- shutdownRequested:
      log.Info("Shutting down client handler")
      return
    }
  }
}

func handleClient(conn net.Conn, c chan *email.SMTPEmail) {
  client := client.MakeClient(conn)
  defer client.Close()

  // close on shutdown using strategy from http://stackoverflow.com/a/13419724
  quit := false
  id, shutdownRequested := shutdown.AddShutdownListener("Client processor")
  defer shutdown.RoutineDone(id)

  go func() {
    <- shutdownRequested
    quit = true
    client.Close()
  }()

  err := client.BeginReceive(c)
  if err != nil {
    if quit {
      log.Info("Shutting down client")
      return
    }
    log.WithFields(log.Fields{
      "error": err.Error(),
    }).Info("error encountered, closing client connection")
  }
}
