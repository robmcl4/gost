package matchers

import (
  "github.com/jhillyerd/go.enmime"
)

type MatchId string

// Matchers can positively identify an email of interest
type Matcher interface {
  GetId()                   MatchId
  Matches(*enmime.MIMEBody) bool
}
