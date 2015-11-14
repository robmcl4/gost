gost
====

[![Build Status](https://travis-ci.org/robmcl4/gost.svg)](https://travis-ci.org/robmcl4/gost) [![Coverage Status](https://coveralls.io/repos/robmcl4/gost/badge.svg?branch=master&service=github)](https://coveralls.io/github/robmcl4/gost?branch=master)

An SMTP server in Go for testing. Currently in Proof-Of-Concept.

Concept
-------

Gost (Go SMTP Test) is meant to facilitate testing automated emails. The end
goal a server that will accept connections via HTTP for clients (probably test
runners) to register their desire to listen for an email. The server then
receives an email via SMTP and notifies the client. The client will then either
receive the content of the email via the RabbitMQ message, or initiate an HTTP
request for the body of the email.

Client TCP Protocol
-------------------

Testing agents use raw TCP sockets to connect to the Gost server, register
pattern matchers, and await matches from the server.

Each message is terminated by the newline character `'\n'`.

### Handshake

1. Client opens a connection to Gost (default port 60510)
2. Server sends `HELLO`
3. Client selects an appropriate version and sends `OLLEH {"version": 3}`
4. If the server supports the version, server sends `OKAY`, otherwise sends
   `UNSUPPORTED` and drops the connection.
5. Client begins main operative procedure

### Main Operative Procedures, version 1

Client commands:

* `PING`, server responds `PONG`
* `MATCH {"to": "foo@bar.com"}`, server responds `MATCHPUT {"match_id": <uuid>}`

Server commands:

* `PING`, client responds `PONG`
* `GOTEMAIL {"match_id": <uuid>, "email": {email...}}`, client responds `OKAY`
* `EXPIRED {"match_id": <uuid>}`, client responds `OKAY`

### Finishing

When the client is finished, sends `QUIT` and the server
drops the connection.
