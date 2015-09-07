package smtp_server

import (
  "net"
  "github.com/robmcl4/gost/smtp_server/email"
  "github.com/robmcl4/gost/smtp_server/client"
  log "github.com/Sirupsen/logrus"
)

func ReceiveEmail(c chan *email.SMTPEmail) error {
  l, err := getServerConnection()
  if err != nil {
    return err
  }
  cxnChan := make(chan net.Conn, 10)
  go handleClients(cxnChan, c)
  err = listenForConnections(l, cxnChan)
  if err != nil {
    return err
  }
  return nil
}

func handleClients(cxnChan chan net.Conn, emChan chan *email.SMTPEmail) {
  for {
    go handleClient(<- cxnChan, emChan)
  }
}

func handleClient(conn net.Conn, c chan *email.SMTPEmail) {
  client := client.MakeClient(conn)
  err := client.BeginReceive(c)
  if err != nil {
    client.Close()
    log.WithFields(log.Fields{
      "error": err.Error(),
    }).Error("client encountered error while receiving messages")
  }
}
