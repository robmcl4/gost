package smtp_server

import (
  "net"
  "bufio"
  "strings"
  "errors"
  "regexp"
)

type client struct {
  conn net.Conn
  in   *bufio.Reader
  out  *bufio.Writer
}

func makeClient(c net.Conn) client {
  return client{ c,
                 bufio.NewReader(c),
                 bufio.NewWriter(c) }
}

// Gets an SMTP command from the client.
// Verb is the 4-letter SMTP verb.
// Extra is any characters appearing after the verb (trimmed)
// On error condition, responds to client with 500 Syntax Error. Never attempts
// to close the connection.
func (c *client) getCommand() (verb string, extra string, err error) {
  line, err := c.readLine()
  if err != nil {
    return
  }
  err = checkCmdSyntax(line)
  if err != nil {
    c.notifySyntaxError()
    return
  }
  verb, extra = splitVerb(line)
  return
}

// Reads from the client until the next newline character.
// Returns the line read after trimming leading/trailing whitespace.
// On transport error, returns the error.
func (c *client) readLine() (string, error) {
  raw, err := c.in.ReadBytes('\n')
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(raw)), nil
}

func (c *client) notifySyntaxError() {
  c.out.WriteString("500 Syntax Error\n")
  c.out.Flush()
}

// -----------------------------------------------------------------------------

// Checks that a command is formatted exactly "ABCD" or "ABCD EFG XYZ"
func checkCmdSyntax(s string) error {
  if !checkCmdSyntaxRegexp.MatchString(s) {
    return errors.New("Syntax Error")
  }
  return nil
}
var checkCmdSyntaxRegexp = regexp.MustCompile(`^[A-Z]{4}( .+$|$)`)

// Splits the command into verb and extra parts.
// Precondition: s is at least 4 characters in length
func splitVerb(s string) (verb string, extra string) {
  verb = s[:4]
  extra = ""
  if len(s) > 5 {
    extra = strings.TrimSpace(s[5:])
  }
  return
}
