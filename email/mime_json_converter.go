package email

import (
  "encoding/json"
  "encoding/base64"
  "github.com/jhillyerd/go.enmime"
)

// Converts the give MIMEBody to JSON.
// This method will *not* include every aspect of the email, only a select
// subset of attributes. So, this should not be used as a serialization method
// with expectation of full support for deserialization.
func dumpMimeToJson(body *enmime.MIMEBody) ([]byte, error) {
  root := make(map[string]interface{})
  root["html"] = body.Html
  root["text"] = body.Text
  root["attachments"] = convertAttachments(body.Attachments)
  root["inlines"] = convertAttachments(body.Inlines)
  return json.Marshal(root)
}

func convertAttachments(attachments []enmime.MIMEPart) [](map[string]interface{}) {
  if attachments == nil {
    return [](map[string]interface{}){}
  }
  ret := make([](map[string]interface{}), 0, len(attachments))
  for _, attch := range attachments {
    ret = append(ret, convertAttachment(attch))
  }
  return ret
}

func convertAttachment(attachment enmime.MIMEPart) map[string]interface{} {
  // sanity check
  if attachment == nil {
    return nil
  }

  ret := make(map[string]interface{})
  if attachment.FileName() != "" {
    ret["fileName"] = attachment.FileName()
  } else {
    ret["fileName"] = nil
  }
  if attachment.ContentType() != "" {
    ret["contentType"] = attachment.ContentType()
  } else {
    ret["contentType"] = nil
  }
  if content := attachment.Content(); content != nil && len(content) > 0 {
    ret["content"] = base64.StdEncoding.EncodeToString(content)
  } else {
    ret["content"] = nil
  }
  return ret
}
