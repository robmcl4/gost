package memory_storage

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestInitializeDoesSomething(t *testing.T) {
  mb := NewMemoryBackend()
  assert.False(t, mb.initialized, "should not be marked as initialized")
  mb.Initialize()
  assert.True(t, mb.initialized, "should be marked as initialized")
  mb.Shutdown()
}
