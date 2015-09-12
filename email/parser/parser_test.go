package parser

import (
  "testing"
  "github.com/robmcl4/gost/email"
  "github.com/stretchr/testify/assert"
)

func TestParseHeadersSingleHeader(t *testing.T) {
  em := getEm("Content-Type: text/plain\r\n\r\nThis is the body.")
  dest := new(email.ParsedEmail)

  assert.Nil(
    t,
    parseHeaders(em, dest),
    "should have no errors",
  )

  assert.Len(
    t,
    dest.Headers,
    1,
  )
  assert.Equal(
    t,
    email.Header{
      "Content-Type",
      " text/plain",
    },
    dest.Headers[0],
  )
}

func TestParseHeadersSingleMultilineHeader(t *testing.T) {
  em := getEm("Foo-Bar: line1\r\n  line2\t \r\n  line3\r\n\r\nThe body.")
  dest := new(email.ParsedEmail)

  assert.Nil(
    t,
    parseHeaders(em, dest),
    "should have no errors",
  )

  assert.Len(
    t,
    dest.Headers,
    1,
  )
  assert.Equal(
    t,
    email.Header{
      "Foo-Bar",
      " line1  line2\t   line3",
    },
    dest.Headers[0],
  )
}

func TestParseHeadersMultipleHeaders(t *testing.T) {
  em := getEm("To: <foo@bar.com>\r\n" +
              "Cc: <ding@dong.com>\r\n" +
              "Thing: other\r\n" +
              "  something\r\n" +
              "\r\n" +
              "The body.")
  dest := new(email.ParsedEmail)

  assert.Nil(
    t,
    parseHeaders(em, dest),
    "should have no errors",
  )
  assert.Equal(
    t,
    []email.Header{
      email.Header{
        "To",
        " <foo@bar.com>",
      },
      email.Header{
        "Cc",
        " <ding@dong.com>",
      },
      email.Header{
        "Thing",
        " other  something",
      },
    },
    dest.Headers,
  )
}

func TestParseHeadersNoHeaders(t *testing.T) {
  em := getEm("this is not a header\r\nneither is this")
  assert.NotNil(
    t,
    parseHeaders(em, new(email.ParsedEmail)),
    "should complain about finding no header values",
  )
  em = getEm("  this is not a header but it doe have whitespace\r\nfoo")
  assert.NotNil(
    t,
    parseHeaders(em, new(email.ParsedEmail)),
    "should complain about finding no header values",
  )
  em = getEm("")
  assert.NotNil(
    t,
    parseHeaders(em, new(email.ParsedEmail)),
    "should complain about empty string",
  )
}

func TestIsHeaderStartPositive(t *testing.T) {
  assert.True(
    t,
    isHeaderStart("FOO: bar"),
    `should be true for "FOO: bar"`,
  )
  assert.True(
    t,
    isHeaderStart("foo: bar"),
    `should be true for "foo: bar"`,
  )
  assert.True(
    t,
    isHeaderStart("FOO:"),
    `should be true for "FOO:"`,
  )
  assert.True(
    t,
    isHeaderStart("FOO-BAR: baz"),
    `should be true for "FOO-BAR: baz"`,
  )
  assert.True(
    t,
    isHeaderStart("FOO1-B3AR1: baz"),
    `should be true for "FOO1-B3AR1: baz"`,
  )
  assert.True(
    t,
    isHeaderStart("0123456789: BAZ"),
    `should be true for "FOO1-B3AR1: baz"`,
  )
  assert.True(
    t,
    isHeaderStart("FOO: BAR\tBAZ"),
    `should be true for "FOO: BAR\tBAZ"`,
  )
}

func TestIsHeaderStartNegative(t *testing.T) {
  assert.False(
    t,
    isHeaderStart("FOO_BAR_BAZ: ding"),
    `should be false for "FOO_BAR_BAZ: ding"`,
  )
  assert.False(
    t,
    isHeaderStart("\t\tWHOO: DING"),
    `should be false for "\t\tWHOO: DING"`,
  )
  assert.False(
    t,
    isHeaderStart("FOO\tBAR: baz"),
    `should be false for "FOO\tBAR: baz"`,
  )
  assert.False(
    t,
    isHeaderStart(""),
    `should be false for empty string`,
  )
}

func TestIsHeaderContinuationPositive(t *testing.T) {
  assert.True(
    t,
    isHeaderContinuation("\tfoo bar baz"),
    `should be true for "\t\tfoo: bar baz"`,
  )
  assert.True(
    t,
    isHeaderContinuation(" bazinga"),
    `should be true for " bazinga"`,
  )
  assert.True(
    t,
    isHeaderContinuation("  \t  yasdfiajsdfijasdifub"),
    `should be true for "  \t  yasdfiajsdfijasdifub"`,
  )
}

func TestIsHeaderContinuationNegative(t *testing.T) {
  assert.False(
    t,
    isHeaderContinuation(""),
    `should be false for empty string`,
  )
  assert.False(
    t,
    isHeaderContinuation("nospaces"),
    `should be false with no leading spaces`,
  )
}

func TestIsHeaderEndPositive(t *testing.T) {
  assert.True(
    t,
    isHeaderEnd(""),
    `is true for empty string`,
  )
}

func TestIsHeaderEndNegative(t *testing.T) {
  assert.False(
    t,
    isHeaderEnd("somecharacters"),
    `is false for a line with word characters`,
  )
  assert.False(
    t,
    isHeaderEnd("\t\tfoobar"),
    `is false for a line starting with whitespace`,
  )
}

func getEm(dat string) *email.SMTPEmail {
  return &email.SMTPEmail{
    "foo@bar.com",
    []string{"bar@example.com"},
    dat,
  }
}
