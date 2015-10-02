// Package registry contains the global collection of email matchers.
// The matcher collection is represented as a linked-list.
package registry

import (
  "time"
  "sync"
  "github.com/jhillyerd/go.enmime"
  "github.com/robmcl4/gost/matchers"
  "github.com/robmcl4/gost/config"
  log "github.com/Sirupsen/logrus"
)

// An element of the linked list of matchers.
type matcherLLElem struct {
  expiry  time.Time
  matcher matchers.Matcher
  next    *matcherLLElem
}

// A Linked-List to hold the global collection of currently
// active matchers. A client is expected to use the list in only
// two manners:
// 1. iterating forward-to-backward
// 2. inserting to the end of the list
// 3. deleting from the front of the list
// Therefore, a linked list can be a reasonably fast
// implementation of a global concurrent collection.
var matcherListHead *matcherLLElem
var matcherListTail *matcherLLElem
var matcherListLock = sync.Mutex{}
var matcherListSize = int64(0)

// Checks for a match of this email among the global matchers.
// Returns an array of the matched Ids.
func GetMatches(e *enmime.MIMEBody) []matchers.MatchId {
  ret := make([]matchers.MatchId, 0)

  // we don't actually need a lock to access the head, so hop right on
  for curr := matcherListHead; curr != nil; curr = curr.next {
    if curr.matcher.Matches(e) {
      ret = append(ret, curr.matcher.GetId())
    }
  }

  return ret
}

// Inserts a matcher.
func InsertMatcher(m matchers.Matcher) {
  newElem := &matcherLLElem{
    expiry:  time.Now().Add(time.Duration(config.GetEmailTTL())*time.Second),
    matcher: m,
    next:    nil,
  }
  matcherListLock.Lock()
  if matcherListTail == nil {
    matcherListHead = newElem
  } else {
    matcherListTail.next = newElem
  }
  matcherListTail = newElem
  matcherListSize++
  matcherListLock.Unlock()

  log.WithFields(log.Fields{
    "matcherId": m.GetId(),
  }).Info("Registered matcher")
}

// Garbage collects expired matchers.
func GarbageCollect() {
  matcherListLock.Lock()

  for curr := matcherListHead;
      curr != nil && curr.expiry.Before(time.Now());
      curr = curr.next {
    matcherListHead = curr.next
    matcherListSize--

    log.WithFields(log.Fields{
      "matcherId": curr.matcher.GetId(),
    }).Info("Retiring matcher")
  }

  // if everything was removed, then set the tail to nil also
  if matcherListHead == nil {
    matcherListTail = nil
  }

  matcherListLock.Unlock()
}

// Gets the number of matchers in the collection.
// This is not thread-safe, and should only be used as an estimate
// in most cases.
func Size() int64 {
  return matcherListSize
}

// Removes all matchers in the collection.
func Clear() {
  matcherListLock.Lock()
  matcherListHead = nil
  matcherListTail = nil
  matcherListSize = 0
  matcherListLock.Unlock()
}
