package client

import (
  "net"
  "bufio"
  "fmt"
  log "github.com/Sirupsen/logrus"
)

type Client struct {
  conn net.Conn
  in   *bufio.Reader
  out  *bufio.Writer
}

type middlewareReadWriter struct {
  conn net.Conn
}

func (mrw *middlewareReadWriter) Read(p []byte) (n int, err error) {
  n, err = mrw.conn.Read(p)
  log.WithFields(log.Fields{
    "bytesRead": n,
    "error": err,
    "bytes": fmt.Sprintf("%s", p[:n]),
  }).Debug("read data from client")
  return
}

func (mrw *middlewareReadWriter) Write(p []byte) (n int, err error) {
  n, err = mrw.conn.Write(p)
  log.WithFields(log.Fields{
    "bytesWritten": n,
    "error": err,
    "bytes": fmt.Sprintf("%s", p[:n]),
  }).Debug("wrote data to client")
  return
}

func MakeClient(c net.Conn) *Client {
  mrw := &middlewareReadWriter{c}
  return &Client{ c,
                  bufio.NewReaderSize(mrw, 512),
                  bufio.NewWriter(mrw) }
}
