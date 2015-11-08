package email_matchers

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/robmcl4/gost/email"
)

func TestPositiveMatch(t *testing.T) {
  e_smtp := email.SMTPEmail{
    Contents: []byte(
`Foo: bar
Bar: baz
The-Thing: Other-Thing
To: The Dude <me@mysite.com>

A body to the message
`),
  }
  e, err := e_smtp.Parse()
  assert.Nil(t, err)
  matcher := NewToMatcher("me@mysite.com")
  assert.True(t, matcher.Matches(e), "should match email")
}

func TestNegativeMatch(t *testing.T) {
  e_smtp := email.SMTPEmail{
    Contents: []byte(
`Foo: bar
Bar: baz
The-Thing: Other-Thing
To: The Dude <me@mysite.com>

A body to the message
`),
  }
  e, err := e_smtp.Parse()
  assert.Nil(t, err)
  matcher := NewToMatcher("NOTme@mysite.com")
  assert.False(t, matcher.Matches(e), "should match email")
}

func TestPositiveMatchMultiple(t *testing.T) {
  e_smtp := email.SMTPEmail{
    Contents: []byte(
`Foo: bar
Bar: baz
The-Thing: Other-Thing
To: "The Dude" <me@mysite.com>, "The Other Dude" <tim@mysite.com>

A body to the message
`),
  }
  e, err := e_smtp.Parse()
  assert.Nil(t, err)
  matcher := NewToMatcher("tim@mysite.com")
  assert.True(t, matcher.Matches(e), "should match email")
}

func TestNegativeMatchMultiple(t *testing.T) {
  e_smtp := email.SMTPEmail{
    Contents: []byte(
`Foo: bar
Bar: baz
The-Thing: Other-Thing
To: "The Dude" <me@mysite.com>, "The Other Dude" <tim@mysite.com>

A body to the message
`),
  }
  e, err := e_smtp.Parse()
  assert.Nil(t, err)
  matcher := NewToMatcher("foo@mysite.com")
  assert.False(t, matcher.Matches(e), "should match email")
}

func TestNegativeMatchNoAddressees(t *testing.T) {
  e_smtp := email.SMTPEmail{
    Contents: []byte(
`Foo: bar
Bar: baz
The-Thing: Other-Thing

A body to the message
`),
  }
  e, err := e_smtp.Parse()
  assert.Nil(t, err)
  matcher := NewToMatcher("foo@mysite.com")
  assert.False(t, matcher.Matches(e), "should match email")
}

func TestCanGetId(t *testing.T) {
  matcher := NewToMatcher("foo@bar.com")
  assert.NotEmpty(t, matcher.GetId())
  assert.Equal(t, matcher.GetId(), matcher.GetId())
}
