package client

import (
  "fmt"
  "strings"
  "errors"
  "regexp"
  "bytes"
  "github.com/robmcl4/gost/config"
)

// Gets an SMTP command from the client.
// Verb is the 4-letter SMTP verb.
// Extra is any characters appearing after the verb (trimmed).
// When the client requests a NOOP, responds 250 OK and continues reading for
// next command without returning the NOOP.
// On error condition, responds to client with 500 Syntax Error. Never attempts
// to close the connection.
func (c *Client) getCommand() (verb string, extra string, err error) {
  for {
    verb, extra, err = c.getSingleCommand()
    if verb == "NOOP" {
      err = c.notifyOk()
      if err != nil {
        return "", "", err
      }
    } else {
      return
    }
  }
}

// Gets a single command. Unlike getCommand, this may return the NOOP verb.
func (c *Client) getSingleCommand() (verb string, extra string, err error) {
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

// Reads the body of en email that immediately follows the DATA command.
func (c *Client) readDataBody() (string, error) {
  buf := bytes.Buffer{}
  for {
    line, err := c.in.ReadBytes('\n')
    if err != nil {
      return "", err
    }
    if bytes.Equal(line, []byte(".\r\n")) {
      break
    }
    buf.Write(line)
  }
  return buf.String(), nil
}

// Reads from the client until the next newline character.
// Returns the line read after trimming leading/trailing whitespace.
// On transport error, returns the error.
func (c *Client) readLine() (string, error) {
  raw, err := c.in.ReadBytes('\n')
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(raw)), nil
}

func (c *Client) notifySyntaxError() {
  c.out.WriteString("500 Syntax Error\n")
  c.out.Flush()
}

func (c *Client) notifyOk() error {
  _, err := c.out.WriteString("250 Ok\n")
  if err != nil {
    return err
  }
  err = c.out.Flush()
  return err
}

func (c *Client) notifyServiceReady() error {
  _, err := c.out.WriteString(fmt.Sprintf("220 %s\n", config.GetFQDN()))
  if err != nil {
    return err
  }
  err = c.out.Flush()
  return err
}

func (c *Client) notifyBadSequence() error {
  _, err := c.out.WriteString("503 Bad Sequence")
  if err != nil {
    return err
  }
  err = c.out.Flush()
  return err
}

// - Utilities -----------------------------------------------------------------

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