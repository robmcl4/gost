package smtp_server

import (
  "net"
  "fmt"
  "github.com/robmcl4/gost/smtp_server/email"
  "github.com/robmcl4/gost/smtp_server/client"
)

func ReceiveEmail(c chan *email.SMTPEmail) error {
  l, err := getServerConnection()
  if err != nil {
    return err
  }
  cxnChan := make(chan net.Conn, 10)
  go listenForConnections(l, cxnChan)
  for {
    conn := <- cxnChan
    go handleClient(conn, c)
  }
}

func handleClient(conn net.Conn, c chan *email.SMTPEmail) {
  client := client.MakeClient(conn)
  err := client.BeginReceive(c)
  if err != nil {
    client.Close()
    fmt.Printf("ERROR: %s\n", err.Error())
  }
}
