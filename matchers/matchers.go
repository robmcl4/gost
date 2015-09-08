package matchers

import (
  "time"
  "sync"
  "github.com/satori/go.uuid"
  "github.com/robmcl4/gost/email"
  "github.com/robmcl4/gost/config"
)

type Matcher interface {
  Matches(*email.SMTPEmail) bool
}

// An element of the linked list of matchers.
type matcherLLElem struct {
  id      string
  expiry  time.Time
  matcher Matcher
  next    *matcherLLElem
}

// A Linked-List to hold the global collection of currently
// active matchers. A client is expected to use the list in only
// two manners:
// 1. iterating forward-to-backward, possibly deleting elements and
// 2. inserting to the end of the list.
// Therefore, a linked list can be a reasonably fast
// implementation of a global concurrent collection
// with non-blocking reads.
var matcherListHead *matcherLLElem
var matcherListInsertLock = sync.Mutex{}

// Checks for a match of this email among the global matchers.
// Returns a list of all matcher Ids that match the email.
func GetMatches(e *email.SMTPEmail) []string {
  ret := make([]string, 0)

  for curr := matcherListHead; curr != nil; curr = curr.next {
    if curr.matcher.Matches(e) {
      ret = append(ret, curr.id)
    }
  }

  return ret
}

// Inserts a matcher, returns the matcher's Id
func InsertMatcher(m Matcher) string {
  newElem := &matcherLLElem{
    uuid.NewV4().String(),
    time.Now().Add(time.Duration(config.GetEmailTTL())*time.Second),
    m,
    nil,
  }
  matcherListInsertLock.Lock()
  newElem.next = matcherListHead
  matcherListHead = newElem
  matcherListInsertLock.Unlock()

  return newElem.id
}
