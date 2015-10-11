package matchers

import (
  "github.com/jhillyerd/go.enmime"
  "github.com/satori/go.uuid"
)

type MatchId string

// Matchers can positively identify an email of interest
type Matcher interface {
  GetId()                   MatchId
  Matches(*enmime.MIMEBody) bool
}

// Base matcher with lazy auto-generation for the ID
type BaseMatcher struct {
  id MatchId
}

func (b *BaseMatcher) GetId() MatchId {
  if b.id == "" {
    b.id = MatchId(uuid.NewV4().String())
  }
  return b.id
}
