package client

import (
  "net"
  "bytes"
  "bufio"
  "testing"
  "errors"
  "github.com/robmcl4/gost/email"
  "github.com/stretchr/testify/assert"
)

type myConn struct {
  net.Conn
  closed bool
}

func (m *myConn) Close() error {
  m.closed = true
  return nil
}

type errConn struct {
  net.Conn
}

func (_ *errConn) Close() error {
  return errors.New("My Close Error")
}

type errWriter struct { }

func (_ *errWriter) Write(_ []byte) (int, error) {
  return 0, errors.New("Write error!")
}

type errReader struct { }

func (_ *errReader) Read(_ []byte) (int, error) {
  return 0, errors.New("Read error!")
}

func TestCloseTerminatesConnection(t *testing.T) {
  mybuf := new(bytes.Buffer)
  mycon := myConn{}
  c := Client{
    &mycon,
    bufio.NewReader(new(bytes.Buffer)),
    bufio.NewWriter(mybuf),
  }
  err := c.Close()
  assert.Equal(
    t,
    "421 Service Unavailable: Terminating Connection\r\n",
    mybuf.String(),
  )
  assert.NoError(t, err)
  assert.True(t, mycon.closed, "Close() should have been called")
}

func TestCloseBubblesUpError(t *testing.T) {
  mybuf := new(bytes.Buffer)
  mycon := errConn{}
  c := Client{
    &mycon,
    bufio.NewReader(new(bytes.Buffer)),
    bufio.NewWriter(mybuf),
  }
  err := c.Close()
  assert.Equal(
    t,
    "421 Service Unavailable: Terminating Connection\r\n",
    mybuf.String(),
  )
  assert.Error(t, err)
  assert.Equal(t, err.Error(), "My Close Error")
}

func TestCloseIgnoresWriteErrors(t *testing.T) {
  mycon := myConn{}
  c := Client{
    &mycon,
    bufio.NewReader(new(bytes.Buffer)),
    bufio.NewWriter(new(errWriter)),
  }
  err := c.Close()
  assert.NoError(t, err)
  assert.True(t, mycon.closed, "Close() should have been called")
}

func TestOpensWithServiceReady(t *testing.T) {
  mybuf := new(bytes.Buffer)
  c := Client{
    &myConn{},
    bufio.NewReader(new(bytes.Buffer)),
    bufio.NewWriter(mybuf),
  }
  ch := make(chan *email.SMTPEmail, 1)
  assert.Error(t, c.BeginReceive(ch))
  assert.Len(t, ch, 0)
  assert.Equal(t, "220 mail.example.com ESMTP\r\n", mybuf.String())
}

func TestPassesErrorIfCannotWrite(t *testing.T) {
  c := Client{
    &myConn{},
    bufio.NewReader(new(bytes.Buffer)),
    bufio.NewWriter(new(errWriter)),
  }
  ch := make(chan *email.SMTPEmail, 1)
  err := c.BeginReceive(ch)
  assert.Error(t, err)
  assert.Equal(t, "Write error!", err.Error())
  assert.Len(t, ch, 0)
}

func TestHelo(t *testing.T) {
  ch, out, err := getClientOutput("HELO mail.example.com\r\n")
  assert.Error(t, err)
  assert.Len(t, ch, 0)
  assert.Equal(
    t,
    "220 mail.example.com ESMTP\r\n" +
    "250 Ok\r\n",
    out,
  )
}

func TestEhlo(t *testing.T) {
  ch, out, err := getClientOutput("EHLO\r\n")
  assert.Error(t, err)
  assert.Len(t, ch, 0)
  assert.Equal(
    t,
    "220 mail.example.com ESMTP\r\n" +
    "250-mail.example.com supports ONE extension:\r\n" +
    "250-8BITMIME\r\n",
    out,
  )
}

func TestUnknownHandshakeVerb(t *testing.T) {
  ch, out, err := getClientOutput("FOOO\r\n")
  assert.Error(t, err)
  assert.Len(t, ch, 0)
  assert.Equal(
    t,
    "220 mail.example.com ESMTP\r\n" +
    "503 Bad Sequence\r\n",
    out,
  )
}

func TestSendSimpleEmail(t *testing.T) {
  ch, out, err := getClientOutput(
    "HELO mail.example.com\r\n" +
    "MAIL FROM:<foo@example.com>\r\n" +
    "RCPT TO:<alice@example.com>\r\n" +
    "DATA\r\n" +
    "From: <foo@example.com>\r\n" +
    "To: <alice@example.com>\r\n" +
    "Subject: Test\r\n" +
    "\r\n" +
    "Hey what's up\r\n" +
    ".\r\n" +
    "QUIT\r\n",
  )
  assert.Error(t, err)
  assert.Equal(t, "Client asked to quit", err.Error())
  assert.Len(t, ch, 1)
  assert.Equal(
    t,
    "220 mail.example.com ESMTP\r\n" +
    "250 Ok\r\n" +
    "250 Ok\r\n" +
    "250 Ok\r\n" +
    "354 Start Mail Input\r\n" +
    "250 Ok\r\n",
    out,
  )
  em := <- ch
  assert.NotNil(t, em)
  assert.Len(t, em.To, 1)
  assert.Equal(t, "alice@example.com", em.To[0])
  assert.Equal(t, "foo@example.com", em.From)
  assert.NotNil(t, em.Contents)
}

func getClientOutput(in string) (ch chan *email.SMTPEmail,
                                 output string,
                                 err error) {
  outb := new(bytes.Buffer)
  inb := bytes.NewBufferString(in)
  c := Client{
    &myConn{},
    bufio.NewReader(inb),
    bufio.NewWriter(outb),
  }
  ch = make(chan *email.SMTPEmail, 10)
  err = c.BeginReceive(ch)
  output = outb.String()
  return
}
