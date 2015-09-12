package parser

import (
  "regexp"
  "strings"
  "errors"
  "fmt"
  "github.com/robmcl4/gost/email"
)

func Parse(e *email.SMTPEmail) (*email.ParsedEmail, error) {
  ret := new(email.ParsedEmail)
  ret.Original = e
  if err := parseHeaders(e, ret); err != nil {
    return nil, err
  }
  if err := parseParts(e, ret); err != nil {
    return nil, err
  }
  return ret, nil
}

func parseHeaders(src *email.SMTPEmail, dst *email.ParsedEmail) error {
  var headers []email.Header
  raw := src.Data

  for {
    parts := strings.SplitN(raw, "\n", 2)
    if parts == nil || len(parts) != 2 {
      return errors.New("Reached end of message, found no end to headers")
    }
    line := strings.Trim(parts[0], "\r\n")
    raw = parts[1]

    if isHeaderStart(line) {
      // ok, let's get the header value!
      headers = append(headers, email.Header{
        mustGetHeaderName(line),
        mustGetHeaderValue(line),
      })
    } else if isHeaderContinuation(line) {
      if len(headers) < 1 {
        return errors.New("Expected header, found line with whitespace.")
      }
      // append this header to the end of what we have already
      headers[len(headers)-1].Val += line
    } else if isHeaderEnd(line) {
      break
    } else {
      return fmt.Errorf("Unrecognized header line: %q", line)
    }
  }

  dst.Headers = headers
  return nil
}

// Returns the value of the given header parameter.
// `s` is assumed to be a string of the format suitable for
// `isHeaderStart()` to return true. Specifically, that the `:`
// character terminates the header name field.
func mustGetHeaderValue(s string) string {
  i := strings.Index(s, ":")
  return s[i+1:]
}

// Returns the name of the given parameter.
// `s` is assumed to be a string of the format suitable for
// `isHeaderStart()` to return true. Specifically, that the `:`
// character terminates the header name field.
func mustGetHeaderName(s string) string {
  i := strings.Index(s, ":")
  return strings.TrimSpace(s[0:i])
}

// Returns True if the given string is a line that ends headers.
// I.E., it is the empty string.
func isHeaderEnd(s string) bool {
  return s == ""
}

// Returns True if the given string could continue a header.
// I.E., starts with spaces or tabs.
func isHeaderContinuation(s string) bool {
  return isHeaderContinuationRegex.Match([]byte(s))
}
var isHeaderContinuationRegex = regexp.MustCompile(`^(\s|\t)`)

// Returns True if the given string starts a header.
// I.E., starts with only alphanumeric characters plus hyphens, followed by
// a colon.
func isHeaderStart(s string) bool {
  return isHeaderStartRegex.Match([]byte(s))
}
var isHeaderStartRegex = regexp.MustCompile(`^[A-Za-z0-9\-]+:.*`)

// Parses the Parts part into the dst from src's Data.
// If this email is MIME-multipart then it parses
func parseParts(src *email.SMTPEmail, dst *email.ParsedEmail) error {
  return nil
}
