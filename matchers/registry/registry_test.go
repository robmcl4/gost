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
