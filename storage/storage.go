// Contains all backends for gost.
// Backends are basic storage engines that can store a received email
// and retreive the email at a later date. When an email is stored
// in a backend it should be guaranteed persisted for at least
// `config.GetEmailTTL()` seconds.
package storage

import (
  "errors"
  log "github.com/Sirupsen/logrus"
  "github.com/robmcl4/gost/storage/memory_storage"
  "github.com/robmcl4/gost/config"
  "github.com/robmcl4/gost/email"
)

// An abstract representation of an email storage backend.
type Backend interface {
  // Puts the given email into the backend, returning its unique ID.
  PutEmail(e *email.SMTPEmail) (id email.EmailId, err error)
  // Gets the given email from the backend by its unique ID.
  // If the email is not found, returns nil and nil.
  // The struct retreived may or may not be the same struct stored,
  // but they must have the same data.
  GetEmail(id email.EmailId) (e *email.SMTPEmail, err error)
  // Initializes the backend.
  // If an error returns, the backend should NOT be used.
  // Initialization should NOT be re-attempted.
  Initialize() (error)
  // Marks the backend for shutdown. Function blocks until shutdown successful.
  Shutdown() (error)
}

// Gets the backend as set in configuration.
// The backend returned is already initialized.
func GetBackend() (Backend, error) {
  btype := config.GetBackendType()
  log.WithFields(log.Fields{
    "backend": btype,
  }).Info("creating backend")

  switch btype {
  case "memory":
    m := memory_storage.NewMemoryBackend()
    err := m.Initialize()
    if err != nil {
      return nil, err
    }
    return m, nil
  }
  return nil, errors.New("could not identify backend type")
}
