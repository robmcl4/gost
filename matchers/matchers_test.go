package matchers

import (
  "sync"
  "time"
  "testing"
  "github.com/robmcl4/gost/email"
  "github.com/stretchr/testify/assert"
  "github.com/robmcl4/gost/config"
)

func TestListHeadStartsNil(t *testing.T) {
  assert.Nil(
    t,
    matcherListHead,
  )
}

func TestGetMatchesEmptyList(t *testing.T) {
  matcherListHead = nil
  assert.Len(
    t,
    GetMatches(new(email.SMTPEmail)),
    0,
  )
}

type positiveMatcher struct { }

func (p *positiveMatcher) Matches(e *email.SMTPEmail) bool {
  return true
}

func TestGetSingleMatch(t *testing.T) {
  matcherListHead = &matcherLLElem{
    "mytotallycoolid",
    time.Now().Add(4*time.Minute),
    new(positiveMatcher),
    nil,
  }
  got := GetMatches(new(email.SMTPEmail))
  assert.Len(t, got, 1)
  assert.Contains(t, got, "mytotallycoolid")
  matcherListHead = nil
}

func TestGetMultipleMatch(t *testing.T) {
  matcherListHead = &matcherLLElem{
    "mycoolid2",
    time.Now().Add(4*time.Minute),
    new(positiveMatcher),
    nil,
  }
  matcherListHead = &matcherLLElem{
    "mycoolid1",
    time.Now().Add(4*time.Minute),
    new(positiveMatcher),
    matcherListHead,
  }
  got := GetMatches(new(email.SMTPEmail))
  assert.Len(t, got, 2)
  assert.Contains(t, got, "mycoolid1")
  assert.Contains(t, got, "mycoolid2")
  matcherListHead = nil
}

type negativeMatcher struct { }

func (n *negativeMatcher) Matches(e *email.SMTPEmail) bool {
  return false
}

func TestGetNoMatch(t *testing.T) {
  matcherListHead = &matcherLLElem{
    "mycoolid1",
    time.Now().Add(4*time.Minute),
    new(negativeMatcher),
    nil,
  }
  got := GetMatches(new(email.SMTPEmail))
  assert.Len(t, got, 0)
  matcherListHead = nil
}

func TestInsertMatcherEmptyList(t *testing.T) {
  matcherListHead = nil
  m   := new(negativeMatcher)
  id  := InsertMatcher(m)
  exp := time.Now().Add(time.Duration(config.GetEmailTTL())*time.Second)

  assert.Len(t, id, 36)
  assert.NotNil(t, matcherListHead)
  assert.Nil(t, matcherListHead.next)
  assert.Equal(t, m, matcherListHead.matcher)
  assert.InDelta(
    t,
    0,
    int64(exp.Sub(matcherListHead.expiry)), float64(10*time.Second),
  )
  matcherListHead = nil
}

func TestInsertMatcherLocks(t *testing.T) {
  wg := new(sync.WaitGroup)
  matcherListHead = nil
  inserter := func() {
    InsertMatcher(new(negativeMatcher))
    wg.Done()
  }
  // insert a bunch of elements in their own goroutines
  for i := 0; i<100; i++ {
    wg.Add(1)
    go inserter()
  }
  // make sure we have 50 elements in the list
  wg.Wait()
  count := 0
  for cur := matcherListHead; cur != nil; cur = cur.next {
    count++
  }
  assert.Equal(t, 100, count)
}
