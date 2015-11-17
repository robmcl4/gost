package client

import (
  "io"
  "bufio"
  "fmt"
  log "github.com/Sirupsen/logrus"
)

const BUFFER_SIZE = 512

type Client struct {
  conn io.ReadWriteCloser
  in   *bufio.Reader
  out  *bufio.Writer
}

type middlewareReadWriter struct {
  conn io.ReadWriteCloser
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

func MakeClient(c io.ReadWriteCloser) *Client {
  mrw := &middlewareReadWriter{c}
  return &Client{ c,
                  bufio.NewReaderSize(mrw, BUFFER_SIZE),
                  bufio.NewWriter(mrw) }
}
