package main

import (
  "fmt"
  "github.com/robmcl4/gost/storage"
  "github.com/robmcl4/gost/email"
  "github.com/robmcl4/gost/smtp_server"
  log "github.com/Sirupsen/logrus"
)

func main() {
  log.SetLevel(log.DebugLevel)
  log.Info("Starting gost server")

  backend, err := storage.GetBackend()
  if err != nil {
    log.WithFields(log.Fields{
      "error": err.Error(),
    }).Error("could not initialize backend")
    return
  }

  emch := make(chan *email.SMTPEmail, 64)
  storage.Intercept(backend, emch)

  err = smtp_server.ReceiveEmail(emch)
  if err != nil {
    log.WithFields(log.Fields{
      "error": err.Error(),
    }).Error("could not start connection")
    return
  }
}

func emailReceiver(c chan *email.SMTPEmail) {
  for {
    eml := <- c
    log.WithFields(log.Fields{
      "email": fmt.Sprintf("%+v", eml),
    }).Info("Got email")
  }
}
