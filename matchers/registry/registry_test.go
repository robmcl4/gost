package registry

import (
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

func TestInsertMatchers(t *testing.T) {
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
  assert.Equal(t, Size(), int64(1000), "should have 1000 things in it")
}
