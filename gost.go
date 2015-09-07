package main

import (
  "fmt"
  "github.com/robmcl4/gost/smtp_server/email"
  "github.com/robmcl4/gost/smtp_server"
  log "github.com/Sirupsen/logrus"
)

func main() {
  log.SetLevel(log.DebugLevel)
  log.Info("Starting gost server")
  c := make(chan *email.SMTPEmail, 10)
  go emailLogger(c)
  err := smtp_server.ReceiveEmail(c)
  if err != nil {
    log.WithFields(log.Fields{
      "error": err.Error(),
    }).Error("could not start connection")
    return
  }
}

func emailLogger(c chan *email.SMTPEmail) {
  for {
    eml := <- c
    log.WithFields(log.Fields{
      "email": fmt.Sprintf("%+v", eml),
    }).Info("Got email")
  }
}
