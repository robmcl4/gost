package email_matchers

import (
  "github.com/robmcl4/gost/matchers"
  "github.com/jhillyerd/go.enmime"
)

type ToMatcher struct {
  matchers.BaseMatcher
  addressee string
}

func NewToMatcher(addressee string) *ToMatcher {
  return &ToMatcher{addressee: addressee}
}

// Returns true if this email is addressed to the matcher's addressee
func (t *ToMatcher) Matches(e *enmime.MIMEBody) bool {
  addressees, err := e.AddressList("to")
  if err != nil || addressees == nil {
    return false
  }

  for _, address := range addressees {
    if address.Address == t.addressee {
      return true
    }
  }
  return false
}
