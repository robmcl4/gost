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
  assert.Equal(t, "500 Syntax Error\r\n", mybuf.String())
}

func TestClientNotifyOk(t *testing.T) {
  mybuf := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(new(bytes.Buffer)), bufio.NewWriter(mybuf)}
  c.notifyOk()
  assert.Equal(t, "250 Ok\r\n", mybuf.String())
}

func TestClientNotifyServiceReady(t *testing.T) {
  mybuf := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(new(bytes.Buffer)), bufio.NewWriter(mybuf)}
  c.notifyServiceReady()
  assert.Equal(t, "220 mail.example.com\r\n", mybuf.String())
}

func TestClientNotifyBadSequence(t *testing.T) {
  mybuf := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(new(bytes.Buffer)), bufio.NewWriter(mybuf)}
  c.notifyBadSequence()
  assert.Equal(t, "503 Bad Sequence\r\n", mybuf.String())
}

func TestClientNotifyStartMailInput(t *testing.T) {
  mybuf := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(new(bytes.Buffer)), bufio.NewWriter(mybuf)}
  c.notifyStartMailInput()
  assert.Equal(t, "354 Start Mail Input\r\n", mybuf.String())
}

func TestClientNotifyTerminatingConnection(t *testing.T) {
  mybuf := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(new(bytes.Buffer)), bufio.NewWriter(mybuf)}
  c.notifyTerminateConnection()
  assert.Equal(t,
               "421 Service Unavailable: Terminating Connection\r\n",
               mybuf.String())
}

func TestGetCommand(t *testing.T) {
  reader := bytes.NewBufferString("MAIL FROM:<foo@bar.com>\r\n")
  c := Client{nil, bufio.NewReader(reader), bufio.NewWriter(new(bytes.Buffer))}
  verb, extra, err := c.getCommand()
  assert.Nil(t, err, "should have no errors")
  assert.Equal(t, "MAIL", verb)
  assert.Equal(t, "FROM:<foo@bar.com>", extra)
}

func TestGetCommandNOOP(t *testing.T) {
  reader := bytes.NewBufferString("NOOP\r\nMAIL FROM:<foo@bar.com>\r\n")
  writer := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(reader), bufio.NewWriter(writer)}
  verb, extra, err := c.getCommand()
  assert.Nil(t, err, "should have no errors")
  assert.Equal(t, "MAIL", verb)
  assert.Equal(t, "FROM:<foo@bar.com>", extra)
  assert.Equal(t, "250 Ok\r\n", writer.String())
}

func TestGetCommandQUIT(t *testing.T) {
  reader := bytes.NewBufferString("QUIT\r\n")
  writer := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(reader), bufio.NewWriter(writer)}
  _, _, err := c.getCommand()
  assert.NotNil(t, err)
}


func TestGetCommandError(t *testing.T) {
  reader := bytes.NewBufferString("FO\r\n")
  writer := new(bytes.Buffer)
  c := Client{nil, bufio.NewReader(reader), bufio.NewWriter(writer)}
  verb, extra, err := c.getCommand()
  assert.Equal(t, "", verb, "should have no verb")
  assert.Equal(t, "", extra, "should have no extra")
  assert.NotNil(t, err, "should have an error")
  assert.Equal(t, "500 Syntax Error\r\n", writer.String())
}

func TestReadDataBody(t *testing.T) {
  reader := bytes.NewBufferString("FOO\r\nBAR\r\n.\r\n")
  c := Client{nil, bufio.NewReader(reader), bufio.NewWriter(new(bytes.Buffer))}
  got, err := c.readDataBody()
  assert.Nil(t, err)
  assert.Equal(t, "FOO\r\nBAR\r\n", got)
}

func TestCheckCmdSyntaxTooShort(t *testing.T) {
  assert.NotNil(t, checkCmdSyntax("FO"))
}

func TestCheckCmdSyntaxNoSpace(t *testing.T) {
  assert.NotNil(t, checkCmdSyntax("ABCDEFG"))
}

func TestCheckCmdSyntaxLowercase(t *testing.T) {
  assert.NotNil(t, checkCmdSyntax("abcd"))
  assert.NotNil(t, checkCmdSyntax("abcd efg"))
}
