package email

import (
  "fmt"
  "bytes"
  "net/mail"
  "github.com/jhillyerd/go.enmime"
)

type EmailId string

type SMTPEmail struct {
  To       []string
  From     string
  Contents []byte
  parsed   *enmime.MIMEBody
}

// Parses this email as MIME.
// Returns the email, or an error if one occurred.
func (e *SMTPEmail) Parse() (*enmime.MIMEBody, error) {
  if e.parsed != nil {
    // it's memoized
    return e.parsed, nil
  }

  msg, err := mail.ReadMessage(bytes.NewReader(e.Contents))
  if err != nil {
    return nil, err
  }

  ret, err := enmime.ParseMIMEBody(msg)
  if err != nil {
    return nil, err
  }

  e.parsed = ret
  return ret, nil
}

// Parses the Contents of the email, and converts it to a JSON structure.
// Returns the bytestring of JSON, or an error if conversion failed.
func (e *SMTPEmail) ParseToJson() ([]byte, error) {
  parsed, err := e.Parse()
  if err != nil {
    return nil, err
  }
  return dumpMimeToJson(parsed)
}

func (e *SMTPEmail) String() string {
  return fmt.Sprintf(
    `{to: %v, from: %s, contents: %v}`,
    e.To,
    e.From,
    string(e.Contents),
  )
}
