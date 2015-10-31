package memory_storage

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/robmcl4/gost/email"
)

func TestInitializeDoesSomething(t *testing.T) {
  mb := NewMemoryBackend()
  assert.False(t, mb.initialized, "should not be marked as initialized")
  mb.Initialize()
  defer mb.Shutdown()
  assert.True(t, mb.initialized, "should be marked as initialized")
}

func TestPutEmailSucceeds(t *testing.T) {
  mb := NewMemoryBackend()
  assert.Nil(t, mb.Initialize(), "initialize should not have errors")
  defer mb.Shutdown()
  id, err := mb.PutEmail(
    &email.SMTPEmail{
      To: []string{"foo@bar.com"},
      From: "tim@bob.com",
      Contents: []byte("Foo: bar\r\n\r\nstuff\r\n"),
    },
  )
  assert.NotNil(t, id)
  assert.Nil(t, err)
  assert.True(t, len(id) > 0, "id not empty")
}

func TestGetEmailSucceeds(t *testing.T) {
  mb := NewMemoryBackend()
  assert.Nil(t, mb.Initialize(), "initialize should not have errors")
  defer mb.Shutdown()
  email := &email.SMTPEmail{
    To: []string{"foo@bar.com"},
    From: "tim@bob.com",
    Contents: []byte("Foo: bar\r\n\r\nstuff\r\n"),
  }
  id, err := mb.PutEmail(email)
  assert.NotNil(t, id)
  assert.Nil(t, err)
  assert.True(t, len(id) > 0, "id is not blank")
  got, err := mb.GetEmail(id)
  assert.Nil(t, err)
  assert.NotNil(t, got)
  assert.Equal(t, got, email)
}

func TestGetUnknownEmailReturnsNils(t *testing.T) {
  mb := NewMemoryBackend()
  assert.Nil(t, mb.Initialize(), "initialize should not have errors")
  defer mb.Shutdown()
  got, err := mb.GetEmail(email.EmailId("im-not-an-email-id"))
  assert.Nil(t, got)
  assert.Nil(t, err)
}
