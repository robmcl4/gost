package shutdown

import (
  "testing"
  "os"
  "github.com/stretchr/testify/assert"
)

func TestAddShutdownListener(t *testing.T) {
  setup()

  id, ch := AddShutdownListener("automated test")
  assert.Equal(t, id, ShutdownId(1))
  assert.NotNil(t, ch)
  assert.Len(t, ch, 0)
}

func TestAddShutdownListenerTwice(t *testing.T) {
  setup()

  id1, ch1 := AddShutdownListener("automated test")
  id2, ch2 := AddShutdownListener("automated test")

  assert.NotEqual(t, id1, id2, "ids should be different")
  assert.NotNil(t, ch1)
  assert.NotNil(t, ch2)
}

func TestRoutineDone(t *testing.T) {
  setup()

  id, _ := AddShutdownListener("automated test")

  assert.Len(t, waitgroup.chans, 1)

  RoutineDone(id)

  assert.Len(t, waitgroup.chans, 0)
}

func TestRoutineDoneUnknownId(t *testing.T) {
  setup()

  assert.Len(t, waitgroup.chans, 0)
  RoutineDone(1112111)
  assert.Len(t, waitgroup.chans, 0)
}

func TestRoutineDoneTwice(t *testing.T) {
  setup()

  id1, _ := AddShutdownListener("automated test")
  id2, _ := AddShutdownListener("automated test")

  assert.NotEqual(t, id1, id2)

  RoutineDone(id1)
  RoutineDone(id2)

  assert.Len(t, waitgroup.chans, 0)
}

func TestShutdownNotifiesChans(t *testing.T) {
  setup()

  id, ch := AddShutdownListener("automated test")

  called := false
  go func() {
    called = <- ch
    RoutineDone(id)
  }()

  Shutdown()

  assert.True(t, called, "should have gotten true from the chan")
}

func TestShutdownOnSigint(t *testing.T) {
  setup()

  id, ch := AddShutdownListener("automated test")
  signalCh := make(chan os.Signal, 0)
  oldNotify := notify
  notify = func(c chan<- os.Signal, _ ...os.Signal) {
    go func() {
      c <- <- signalCh
    }()
  }

  go ShutdownOnSigint()
  signalCh <- os.Interrupt

  <- ch
  RoutineDone(id)
  notify = oldNotify
}

func setup() {
  waitgroup = tsWaitGroup{}
  waitgroup.chans = make(map[ShutdownId]chan bool)
}
