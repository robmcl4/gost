package email

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestToJson(t *testing.T) {
  e := new(SMTPEmail)
  e.From = "foo@bar.com"
  e.To = []string{"rcpt1@foo.com", "rcpt2@foo.com"}
  e.Data = "whoaimthedata"
  js := e.ToJson()
  assert.Equal( t,
    `{"data":"whoaimthedata",` +
    `"from":"foo@bar.com",` +
    `"to":["rcpt1@foo.com","rcpt2@foo.com"]}`,
    js,
  )
}
