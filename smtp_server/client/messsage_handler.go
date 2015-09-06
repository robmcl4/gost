package client

import (
  "fmt"
  "regexp"
  "github.com/robmcl4/gost/smtp_server/email"
)

func (c *Client) Close() error {
  c.out.Flush()
  c.notifyTerminateConnection()
  return c.conn.Close()
}

// Begins receiving messages on the client connection.
// Clients are expected to begin with a HELO message.
// Any emails received are put onto the channel.
// On error, returns error and does not attempt to close the connection.
func (c *Client) BeginReceive(ch chan *email.SMTPEmail) error {
  err := c.handleHandshake()
  if err != nil {
    return err
  }
  for {
    eml, err := c.getEmail()
    if err != nil {
      return err
    }
    ch <- eml
  }
}


func (c *Client) handleHandshake() error {
  // send 220 Service Ready
  err := c.notifyServiceReady()
  if err != nil {
    return err
  }
  // Wait for HELO (or in future, EHLO)
  verb, _, err := c.getCommand()
  if err != nil {
    return err
  }
  // If verb is not HELO respond bad sequence
  if verb != "HELO" {
    c.notifyBadSequence()
    return fmt.Errorf("Expected HELO but got %s", verb)
  }
  err = c.notifyOk()
  return err
}


func (c *Client) getEmail() (*email.SMTPEmail, error) {
  ret := new(email.SMTPEmail)
  // get the MAIL command
  verb, extra, err := c.getCommand()
  if verb != "MAIL" {
    c.notifyBadSequence()
    return nil, fmt.Errorf("Expected MAIL but got %s", verb)
  }
  if match := fromRegexp.FindStringSubmatch(extra); match != nil {
    ret.From = match[1]
  }
  err = c.notifyOk()
  if err != nil {
    return nil, err
  }
  // get the RCPT commands
  for {
    verb, extra, err = c.getCommand()
    if verb != "RCPT" {
      break
    }
    if err != nil {
      return nil, err
    }
    if match := toRegexp.FindStringSubmatch(extra); match != nil {
      ret.To = append(ret.To, match[1])
      err = c.notifyOk()
      if err != nil {
        return nil, err
      }
    } else {
      c.notifySyntaxError()
      return nil, fmt.Errorf("Couldn't find recipient email: %s", extra)
    }
  }
  // ok so we should have at least 1 recipient now.. if not, this is an error
  if ret.To == nil || len(ret.To) == 0 {
    c.notifyBadSequence()
    return nil, fmt.Errorf("Expected RCPT command, got %s", verb)
  }
  // this should be the DATA command
  if verb != "DATA" {
    c.notifyBadSequence()
    return nil, fmt.Errorf("Expected DATA command, got %s", verb)
  }
  err = c.notifyStartMailInput()
  if err != nil {
    return nil, err
  }
  // let's receive data...
  ret.Data, err = c.readDataBody()
  if err != nil {
    return nil, err
  }
  err = c.notifyOk()
  if err != nil {
    return nil, err
  }
  return ret, nil
}
var fromRegexp = regexp.MustCompile(`[Ff][Rr][Oo][Mm]:<(.+)>`)
var toRegexp = regexp.MustCompile(`[Tt][Oo]:<(.+)>`)
