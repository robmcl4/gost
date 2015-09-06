package main

import (
  "fmt"
  "github.com/robmcl4/gost/smtp_server/email"
  "github.com/robmcl4/gost/smtp_server"
)

func main() {
  c := make(chan *email.SMTPEmail, 10)
  go emailLogger(c)
  err := smtp_server.ReceiveEmail(c)
  if err != nil {
    fmt.Printf("ERROR starting connection: %s\n", err.Error())
    return
  }
}

func emailLogger(c chan *email.SMTPEmail) {
  for {
    eml := <- c
    fmt.Println("Got email")
    fmt.Printf("%+v\n", eml)
  }
}
