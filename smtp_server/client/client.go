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

type loggedReadWriter struct {
  conn net.Conn
}

func (lrw *loggedReadWriter) Read(p []byte) (n int, err error) {
  n, err = lrw.conn.Read(p)
  log.WithFields(log.Fields{
    "bytesRead": n,
    "error": err,
    "bytes": fmt.Sprintf("%s", p[:n]),
  }).Debug("read data from client")
  return
}

func (lrw *loggedReadWriter) Write(p []byte) (n int, err error) {
  n, err = lrw.conn.Write(p)
  log.WithFields(log.Fields{
    "bytesWritten": n,
    "error": err,
    "bytes": fmt.Sprintf("%s", p[:n]),
  }).Debug("wrote data to client")
  return
}

func MakeClient(c net.Conn) *Client {
  lrw := &loggedReadWriter{c}
  return &Client{ c,
                  bufio.NewReaderSize(lrw, 512),
                  bufio.NewWriter(lrw) }
}
