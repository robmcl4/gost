package memory_storage

import (
  "errors"
  "fmt"
  "sync"
  "time"
  "github.com/satori/go.uuid"
  log "github.com/Sirupsen/logrus"
  "github.com/robmcl4/gost/email"
  "github.com/robmcl4/gost/config"
  "github.com/robmcl4/gost/config/shutdown"
)

// An in-memory store for emails.
type MemoryBackend struct {
  // Read-Write lock for modifying any contents of the backend.
  rwlock       sync.RWMutex
  // Map from Id to Email
  email_by_id  map[email.EmailId]*email.SMTPEmail
  // A queue of email ids and their timestamps, soonest to expire first.
  expiry_queue []timestampedEmail
  // Has this been initialized?
  initialized  bool
  // Channel for manual (non system-wide) shutdown
  quit         chan bool
}

type timestampedEmail struct {
  // The unique Id in the memory backend.
  id         email.EmailId
  // The time this Id should expire.
  expiration time.Time
}

func NewMemoryBackend() *MemoryBackend {
  ret := new(MemoryBackend)
  ret.email_by_id = make(map[email.EmailId]*email.SMTPEmail)
  ret.expiry_queue = make([]timestampedEmail, 0)
  ret.quit = make(chan bool, 0)
  return ret
}

func (b *MemoryBackend) PutEmail(e *email.SMTPEmail) (id email.EmailId, err error) {
  if e == nil {
    return "", errors.New("Cannot put nil email")
  }

  id = email.EmailId(uuid.NewV4().String())
  tsEm := timestampedEmail{id, time.Now().Add(config.GetEmailTTL())}

  // critical section
  b.rwlock.Lock()
  b.email_by_id[id] = e
  b.expiry_queue = append(b.expiry_queue, tsEm)
  b.rwlock.Unlock()

  log.WithFields(log.Fields{
    "email": fmt.Sprintf("%v", e),
    "id": id,
  }).Info("Stored email in memory")

  return
}

func (b *MemoryBackend) GetEmail(id email.EmailId) (e *email.SMTPEmail, err error) {
  b.rwlock.RLock()
  e = b.email_by_id[id]
  b.rwlock.RUnlock()
  log.WithFields(log.Fields{
    "email": fmt.Sprintf("%+v", e),
    "id": id,
  }).Info("Retrieved email")
  return
}

func (b *MemoryBackend) Initialize() error {
  b.rwlock.Lock()
  if b.initialized {
    return nil
  }
  b.initialized = true
  go b.continuousCleanup()
  b.rwlock.Unlock()
  return nil
}

func (b *MemoryBackend) Shutdown() {
  b.quit <- true
}

func (b *MemoryBackend) continuousCleanup() {
  id, shutdownRequested := shutdown.AddShutdownListener("MemoryBackend cleanup")
  defer shutdown.RoutineDone(id)

  for {
    select {
    case <- time.After(config.GetEmailTTL()/4):
      b.doCleanup()
    case <- shutdownRequested:
      log.Info("Shutting down memory cleanup with system")
      return
    case <- b.quit:
      log.Info("Quitting memory cleanup")
      return
    }
  }
}

// Iterates through the expiration queue and evicts expired entries.
// Locks the backend.
func (b *MemoryBackend) doCleanup() {
  log.Info("Cleaning up memory backend")

  b.rwlock.Lock()
  defer b.rwlock.Unlock()

  for len(b.expiry_queue) > 0 && b.expiry_queue[0].expiration.Before(time.Now()) {

    log.WithFields(log.Fields{
      "id": b.expiry_queue[0].id,
    }).Info("Expiring email from memory")

    delete(b.email_by_id, b.expiry_queue[0].id)
    b.expiry_queue = b.expiry_queue[1:]
  }
}
