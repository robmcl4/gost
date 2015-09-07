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
  s, _ := json.Marshal(map[string]interface{}{
    "from": e.From,
    "to": e.To,
    "data": e.Data,
  })
  return string(s)
}
