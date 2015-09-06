package client

import (
  "testing"
  "bytes"
  "bufio"
  "github.com/stretchr/testify/assert"
)

func TestSplitVerbLength4(t *testing.T) {
  verb, extra := splitVerb("ABCD")
  assert.Equal(t, "ABCD", verb)
  assert.Equal(t, "", extra)
}

func TestSplitVerbLength5(t *testing.T) {
  verb, extra := splitVerb("ABCD ")
  assert.Equal(t, "ABCD", verb)
  assert.Equal(t, "", extra)
}

func TestSplitVerbLength6(t *testing.T) {
  verb, extra := splitVerb("ABCD E")
  assert.Equal(t, "ABCD", verb)
  assert.Equal(t, "E", extra)
}

func TestSplitVerbLengthExtraSpace(t *testing.T) {
  verb, extra := splitVerb("ABCD  E")
  assert.Equal(t, "ABCD", verb)
  assert.Equal(t, "E", extra)
}

func TestClientNotifySyntaxError(t *testing.T) {
  mybuf := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(new(bytes.Buffer)), bufio.NewWriter(mybuf)}
  c.notifySyntaxError()
  assert.Equal(t, "500 Syntax Error\n", mybuf.String())
}

func TestClientNotifyOk(t *testing.T) {
  mybuf := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(new(bytes.Buffer)), bufio.NewWriter(mybuf)}
  c.notifyOk()
  assert.Equal(t, "250 Ok\n", mybuf.String())
}

func TestGetCommand(t *testing.T) {
  reader := bytes.NewBufferString("MAIL FROM:<foo@bar.com>\n")
  c := Client{nil, bufio.NewReader(reader), bufio.NewWriter(new(bytes.Buffer))}
  verb, extra, err := c.getCommand()
  assert.Nil(t, err, "should have no errors")
  assert.Equal(t, "MAIL", verb)
  assert.Equal(t, "FROM:<foo@bar.com>", extra)
}

func TestGetCommandNOOP(t *testing.T) {
  reader := bytes.NewBufferString("NOOP\nMAIL FROM:<foo@bar.com>\n")
  writer := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(reader), bufio.NewWriter(writer)}
  verb, extra, err := c.getCommand()
  assert.Nil(t, err, "should have no errors")
  assert.Equal(t, "MAIL", verb)
  assert.Equal(t, "FROM:<foo@bar.com>", extra)
  assert.Equal(t, "250 Ok\n", writer.String())
}

func TestGetCommandError(t *testing.T) {
  reader := bytes.NewBufferString("FO\n")
  writer := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(reader), bufio.NewWriter(writer)}
  verb, extra, err := c.getCommand()
  assert.Equal(t, "", verb, "should have no verb")
  assert.Equal(t, "", extra, "should have no extra")
  assert.NotNil(t, err, "should have an error")
  assert.Equal(t, "500 Syntax Error\n", writer.String())
}