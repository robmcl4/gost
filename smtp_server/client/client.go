package client

import (
  "net"
  "bufio"
)

type Client struct {
  conn net.Conn
  in   *bufio.Reader
  out  *bufio.Writer
}

func MakeClient(c net.Conn) *Client {
  return &Client{ c,
                  bufio.NewReader(c),
                  bufio.NewWriter(c) }
}
