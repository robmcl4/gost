package storage

input (
  "github.com/robmcl4/gost/email"
)

type Backend interface {
  func (b *Backend) PutEmail(e *email.SMTPEmail) (id string, err error)
  func (b *Backend) GetEmail(id string) (e *email.SMTPEmail, err error)
}
