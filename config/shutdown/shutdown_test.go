package shutdown

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestAddShutdownListener(t *testing.T) {
  waitgroup.nextId = 0
  waitgroup.chans = make(map[ShutdownId]chan bool)

  id, ch := AddShutdownListener()
  assert.Equal(t, id, ShutdownId(1))
  assert.NotNil(t, ch)
  assert.Len(t, ch, 0)
}

func TestAddShutdownListenerTwice(t *testing.T) {
  waitgroup.chans = make(map[ShutdownId]chan bool)
  id1, ch1 := AddShutdownListener()
  id2, ch2 := AddShutdownListener()

  assert.NotEqual(t, id1, id2, "ids should be different")
  assert.NotNil(t, ch1)
  assert.NotNil(t, ch2)
}

func TestRoutineDone(t *testing.T) {
  waitgroup.chans = make(map[ShutdownId]chan bool)
  id, _ := AddShutdownListener()

  assert.Len(t, waitgroup.chans, 1)

  RoutineDone(id)

  assert.Len(t, waitgroup.chans, 0)
}

func TestRoutineDoneTwice(t *testing.T) {
  waitgroup.chans = make(map[ShutdownId]chan bool)
  id1, _ := AddShutdownListener()
  id2, _ := AddShutdownListener()

  assert.NotEqual(t, id1, id2)

  RoutineDone(id1)
  RoutineDone(id2)

  assert.Len(t, waitgroup.chans, 0)
}

func TestShutdownNotifiesChans(t *testing.T) {
  waitgroup = tsWaitGroup{}
  waitgroup.chans = make(map[ShutdownId]chan bool)
  id, ch := AddShutdownListener()

  called := false
  cb := func() {
    called = <- ch
    RoutineDone(id)
  }

  go cb()
  Shutdown()

  assert.True(t, called, "should have gotten true from the chan")
}
