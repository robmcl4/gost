gost
====

An SMTP server in Go for testing. Currently in Proof-Of-Concept.

Concept
-------

Gost (Go SMTP Test) is meant to facilitate testing automated emails. The end
goal a server that will accept connections via HTTP for clients (probably test
runners) to register their desire to listen for an email. The server then
receives an email via SMTP and publishes a message, probably on
[RabbitMQ](https://www.rabbitmq.com/). The client will then either receive the
content of the email via the RabbitMQ message, or initiate an HTTP request for
the body of the email.