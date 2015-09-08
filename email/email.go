package email

import (
  "encoding/json"
)

type SMTPEmail struct {
  From string
  To   []string
  Data string
}

func (e *SMTPEmail) ToJson() string {
  s, _ := json.Marshal(e)
  return string(s)
}

func FromJson(s string) (*SMTPEmail, error) {
  ret := new(SMTPEmail)
  err := json.Unmarshal([]byte(s), ret)
  if err != nil {
    return nil, err
  }
  return ret, nil
}