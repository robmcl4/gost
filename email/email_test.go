package email

import (
  "testing"
  "encoding/json"
  "github.com/stretchr/testify/assert"
)

var TEXT_ONLY_EMAIL = []byte("" +
"X-Mozilla-Status: 0001\r\n" +
"X-Mozilla-Status2: 00800000\r\n" +
"X-Mozilla-Keys:    \r\n" +
"To: barbaz@example.com\r\n" +
"From: Foo Bar <foobar@example.com>\r\n" +
"Subject: hello, world\r\n" +
"Message-ID: <56363197.5050006@example.com>\r\n" +
"Date: Sun, 1 Nov 2015 10:36:55 -0500\r\n" +
"User-Agent: Mozilla/5.0 (Windows NT 6.3; WOW64; rv:38.0) Gecko/20100101\r\n" +
" Thunderbird/38.2.0\r\n" +
"MIME-Version: 1.0\r\n" +
"Content-Type: text/plain; charset=utf-8; format=flowed\r\n" +
"Content-Transfer-Encoding: 7bit\r\n" +
"\r\n" +
"hello world\r\n" +
"\r\n" +
"this is a plain text email\r\n" +
"\r\n" +
"- Me\r\n",
)

var TEXT_AND_HTML_EMAIL = []byte("" +
"X-Mozilla-Status: 0001\r\n" +
"X-Mozilla-Status2: 00800000\r\n" +
"X-Mozilla-Keys: \r\n" +
"To: barbaz@example.com\r\n" +
"From: Foo Bar <foobar@example.com>\r\n" +
"Subject: hello, world\r\n" +
"Message-ID: <56363935.6030501@example.com>\r\n" +
"Date: Sun, 1 Nov 2015 11:09:25 -0500\r\n" +
"User-Agent: Mozilla/5.0 (Windows NT 6.3; WOW64; rv:38.0) Gecko/20100101\r\n" +
" Thunderbird/38.2.0\r\n" +
"MIME-Version: 1.0\r\n" +
"Content-Type: multipart/alternative;\r\n" +
" boundary=\"------------070908000008020009090300\"\r\n" +
"\r\n" +
"This is a multi-part message in MIME format.\r\n" +
"--------------070908000008020009090300\r\n" +
"Content-Type: text/plain; charset=utf-8; format=flowed\r\n" +
"Content-Transfer-Encoding: 7bit\r\n" +
"\r\n" +
"this is a message\r\n" +
"\r\n" +
"that should be both html and text\r\n" +
"\r\n" +
"-Me\r\n" +
"\r\n" +
"--------------070908000008020009090300\r\n" +
"Content-Type: text/html; charset=utf-8\r\n" +
"Content-Transfer-Encoding: 7bit\r\n" +
"\r\n" +
"<html>\r\n" +
"  <head>\r\n" +
"\r\n" +
"    <meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\">\r\n" +
"  </head>\r\n" +
"  <body bgcolor=\"#FFFFFF\" text=\"#000000\">\r\n" +
"    this is a message<br>\r\n" +
"    <br>\r\n" +
"    that should be both html and text<br>\r\n" +
"    <br>\r\n" +
"    -Me<br>\r\n" +
"  </body>\r\n" +
"</html>\r\n" +
"\r\n" +
"--------------070908000008020009090300--\r\n",
)

var TEXT_EMAIL_WITH_ATTACHMENT = []byte("" +
"X-Mozilla-Status: 0001\r\n" +
"X-Mozilla-Status2: 00800000\r\n" +
"X-Mozilla-Keys:\r\n" +
"To: barbaz@example.com\r\n" +
"From: Foo Bar <foobar@example.com>\r\n" +
"Subject: sample attachment\r\n" +
"Message-ID: <56363DAC.2030907@example.com>\r\n" +
"Date: Sun, 1 Nov 2015 11:28:28 -0500\r\n" +
"User-Agent: Mozilla/5.0 (Windows NT 6.3; WOW64; rv:38.0) Gecko/20100101\r\n" +
" Thunderbird/38.2.0\r\n" +
"MIME-Version: 1.0\r\n" +
"Content-Type: multipart/mixed;\r\n" +
" boundary=\"------------070600020903010200090505\"\r\n" +
"\r\n" +
"This is a multi-part message in MIME format.\r\n" +
"--------------070600020903010200090505\r\n" +
"Content-Type: text/plain; charset=utf-8; format=flowed\r\n" +
"Content-Transfer-Encoding: 7bit\r\n" +
"\r\n" +
"hello world, I have an attachment\r\n" +
"\r\n" +
"in this email\r\n" +
"\r\n" +
"-me\r\n" +
"\r\n" +
"--------------070600020903010200090505\r\n" +
"Content-Type: image/jpeg;\r\n" +
" name=\"small_image.jpg\"\r\n" +
"Content-Transfer-Encoding: base64\r\n" +
"Content-Disposition: attachment;\r\n" +
" filename=\"small_image.jpg\"\r\n" +
"\r\n" +
"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\r\n" +
"p22oAvvFB/lkx9/2X/MP6wJ3KZ6oZ//Z\r\n" +
"--------------070600020903010200090505--\r\n",
)

func TestCanParsePlainText(t *testing.T) {
  em := &SMTPEmail{Contents: TEXT_ONLY_EMAIL}
  got, err := em.Parse()
  assert.Nil(t, err)
  assert.NotNil(t, got)
  assert.NotEmpty(t, got.Text)
  assert.Empty(t, got.Html)
  assert.Empty(t, got.Attachments)
  assert.Empty(t, got.Inlines)
}

func TestCanParsePlainTextToJson(t *testing.T) {
  em := &SMTPEmail{Contents: TEXT_ONLY_EMAIL}
  got, err := em.ParseToJson()
  assert.Nil(t, err)
  assert.NotNil(t, got)
  // parse json and check its contents
  parsed := make(map[string]interface{})
  assert.Nil(t, json.Unmarshal(got, &parsed), "should unmarshal without errors")
  assert.Empty(t, parsed["html"])
  assert.NotEmpty(t, parsed["text"])
  assert.NotNil(t, parsed["attachments"])
  assert.Len(t, parsed["attachments"], 0)
  assert.NotNil(t, parsed["inlines"])
  assert.Len(t, parsed["inlines"], 0)
}

func TestCanParseHtmlAndText(t *testing.T) {
  em := &SMTPEmail{Contents: TEXT_AND_HTML_EMAIL}
  got, err := em.ParseToJson()
  assert.Nil(t, err)
  assert.NotNil(t, got)
  parsed := make(map[string]interface{})
  assert.Nil(t, json.Unmarshal(got, &parsed), "should unmarshal without errors")
  assert.NotEmpty(t, parsed["html"])
  assert.NotEmpty(t, parsed["text"])
  assert.NotNil(t, parsed["attachments"])
  assert.Len(t, parsed["attachments"], 0)
  assert.NotNil(t, parsed["inlines"])
  assert.Len(t, parsed["inlines"], 0)
}

func TestCanParseTextWithAttachment(t *testing.T) {
  em := &SMTPEmail{Contents: TEXT_EMAIL_WITH_ATTACHMENT}
  got, err := em.ParseToJson()
  assert.Nil(t, err)
  assert.NotNil(t, got)
  parsed := make(map[string]interface{})
  assert.Nil(t, json.Unmarshal(got, &parsed), "should unmarshal without errors")
  assert.Empty(t, parsed["html"])
  assert.NotEmpty(t, parsed["text"])
  assert.NotNil(t, parsed["attachments"])
  assert.Len(t, parsed["attachments"], 1)
  assert.NotNil(t, parsed["inlines"])
  assert.Len(t, parsed["inlines"], 0)
}
