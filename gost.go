package main

import (
  "github.com/robmcl4/gost/storage"
  "github.com/robmcl4/gost/config/shutdown"
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

  go intercept(backend, emch)
  go runServer(emch)
  shutdown.ShutdownOnSigint()
}

func runServer(emch chan *email.SMTPEmail) {
  err := smtp_server.ReceiveEmail(emch)
  if err != nil {
    log.WithFields(log.Fields{
      "error": err.Error(),
    }).Error("could not start connection")
    return
  }
}

func intercept(b storage.Backend, ch chan *email.SMTPEmail) {
  for {
    b.PutEmail(<- ch)
  }
}
