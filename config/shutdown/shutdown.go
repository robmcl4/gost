// Package shutdown contains data and methods to support safe shutdown.
package shutdown

import (
  "sync"
  "sync/atomic"
  "os"
  "os/signal"
  "syscall"
  log "github.com/Sirupsen/logrus"
)

// A unique ID to represent a process' shutdown listener request.
type ShutdownId uint64

// internal type for keeping a threadsafe waitgroup / shutdown requester list
type tsWaitGroup struct {
  sync.WaitGroup
  sync.Mutex
  nextId uint64
  chans  map[ShutdownId]chan bool
}

var waitgroup = tsWaitGroup{chans: make(map[ShutdownId]chan bool)}

// Records that a module will be waiting for shutdown.
// Returns A channel to listen for shutdown requests.
func AddShutdownListener(routineName string) (ShutdownId, chan bool) {
  waitgroup.Lock()
  defer waitgroup.Unlock()

  waitgroup.Add(1)
  ch := make(chan bool, 1)
  id := ShutdownId(atomic.AddUint64(&waitgroup.nextId, 1))
  waitgroup.chans[id] = ch
  log.WithFields(log.Fields{
    "ShutdownId": id,
    "routineName": routineName,
  }).Debug("Added subprocess shutdown listener")
  return id, ch
}

// Marks that a routine has finished shutdown
func RoutineDone(id ShutdownId) {
  waitgroup.Lock()
  defer waitgroup.Unlock()

  if _, hasId := waitgroup.chans[id]; !hasId {
    return
  }

  delete(waitgroup.chans, id)
  log.WithFields(log.Fields{
    "ShutdownId": id,
  }).Debug("Subprocess marked as done")
  waitgroup.Done()
}

// Requests shutdown, and waits until shutdown has completed
func Shutdown() {
  // notify everyone to turn off
  waitgroup.Lock()
  for key := range waitgroup.chans {
    waitgroup.chans[key] <- true
  }
  waitgroup.Unlock()

  waitgroup.Wait()
}

// useful for unit testing so actual OS signals can be swapped out
var notify = signal.Notify

// Blocks and requests shutdown on receiving SIGINT, unblocks on completion
// of shutdown.
func ShutdownOnSigint() {
  ch := make(chan os.Signal, 1)
  notify(ch, os.Interrupt, syscall.SIGTERM)
  <- ch
  log.Info("SIGINT received, shutting down")
  Shutdown()
}
