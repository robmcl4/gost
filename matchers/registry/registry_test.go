package registry

import (
  "time"
  "sync"
  "testing"
  "github.com/jhillyerd/go.enmime"
  "github.com/stretchr/testify/assert"
  "github.com/robmcl4/gost/matchers"
)

type positiveMatcher struct { }

func (p *positiveMatcher) GetId() matchers.MatchId {
  return "1"
}

func (p *positiveMatcher) Matches(e *enmime.MIMEBody) bool {
  return true
}

type negativeMatcher struct  { }

func (n *negativeMatcher) GetId() matchers.MatchId {
  return "2"
}

func (n *negativeMatcher) Matches(e *enmime.MIMEBody) bool {
  return false
}

func TestInsertMatchers(t *testing.T) {
  Clear()

  wg := new(sync.WaitGroup)

  // try to put a bunch of matchers in at once, make sure they
  // all get in
  var putMatcher = func() {
    InsertMatcher(new(positiveMatcher))
    wg.Done()
  }

  for i := 0; i < 1000; i++ {
    wg.Add(1)
    go putMatcher()
  }

  // make sure they all finish
  wg.Wait()

  // make sure the list has 100 things in it
  assert.Equal(t, int64(1000), Size(), "should have 1000 things in it")
}

func TestGetMatchPositive(t *testing.T) {
  Clear()

  InsertMatcher(new(positiveMatcher))

  assert.Equal(t, matchers.MatchId("1"), GetMatches(new(enmime.MIMEBody))[0])
}

func TestGetMatchNegative(t *testing.T) {
  Clear()

  InsertMatcher(new(negativeMatcher))

  assert.Len(t, GetMatches(new(enmime.MIMEBody)), 0)
}

func TestGetPositiveMatches(t *testing.T) {
  Clear()

  InsertMatcher(new(positiveMatcher))
  InsertMatcher(new(positiveMatcher))

  assert.Len(t, GetMatches(new(enmime.MIMEBody)), 2)
}

func TestClearRemovesAll(t *testing.T) {
  Clear()

  assert.Nil(t, matcherListHead)
  assert.Nil(t, matcherListTail)
  assert.Equal(t, matcherListSize, int64(0))

  InsertMatcher(new(positiveMatcher))
  InsertMatcher(new(positiveMatcher))

  Clear()

  assert.Nil(t, matcherListHead)
  assert.Nil(t, matcherListTail)
  assert.Equal(t, matcherListSize, int64(0))
}

func TestGarbageCollectEmpty(t *testing.T) {
  Clear()
  GarbageCollect()
  assert.Nil(t, matcherListHead)
  assert.Nil(t, matcherListTail)
  assert.Equal(t, matcherListSize, int64(0))
}

func TestGarbageCollectsFirstFew(t *testing.T) {
  Clear()
  InsertMatcher(new(positiveMatcher))
  InsertMatcher(new(positiveMatcher))

  matcherListHead.expiry = time.Now().Add(-1000*time.Second)
  GarbageCollect()

  assert.NotNil(t, matcherListHead)
  assert.NotNil(t, matcherListTail)
  assert.Nil(t, matcherListHead.next)
  assert.Nil(t, matcherListTail.next)
  assert.Equal(t, int64(1), matcherListSize)
  assert.Equal(t, matcherListHead, matcherListTail)
}

func TestGarbageCollectsAll(t *testing.T) {
  Clear()
  InsertMatcher(new(positiveMatcher))
  InsertMatcher(new(positiveMatcher))
  matcherListHead.expiry = time.Now().Add(-1000*time.Second)
  matcherListHead.next.expiry = time.Now().Add(-1000*time.Second)

  GarbageCollect()

  assert.Nil(t, matcherListHead)
  assert.Nil(t, matcherListTail)
  assert.EqualValues(t, 0, matcherListSize)
}
