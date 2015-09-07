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
    `{"From":"foo@bar.com",` +
    `"To":["rcpt1@foo.com","rcpt2@foo.com"],` +
    `"Data":"whoaimthedata"}`,
    js,
  )
}

func TestFromJson(t *testing.T) {
  s := `{"From":"foo@bar.com",` +
       `"To":["rcpt1@foo.com","rcpt2@foo.com"],` +
       `"Data":"whoaimthedata"}`
  got, err := FromJson(s)
  assert.Nil(t, err, "should have no errors")
  assert.Equal(t, "foo@bar.com", got.From)
  assert.Equal(t, "rcpt1@foo.com", got.To[0])
  assert.Equal(t, "rcpt2@foo.com", got.To[1])
  assert.Equal(t, "whoaimthedata", got.Data)
}
