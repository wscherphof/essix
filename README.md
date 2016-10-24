# Essix
Package essix runs an essential simple secure stable scalable stateless server.

## Quickstart

1. Install Essix: `go get -u github.com/wscherphof/essix`
1. Initialise a new Essix app: `essix init github.com/you/yourapp`
1. Create a self-signed TLS certificate: `cd $GOPATH/github.com/you/yourapp`, then: `essix cert dev.appsite.com`
1. Create a local one-node Docker swarm: `essix nodes -H dev.appsite.com -m 1 create dev`
1. Install RethinkDB on the swarm: `essix r create dev`
1. Run your app on the swarm: `cd $GOPATH/github.com/you/yourapp`, then: `essix -e DOMAIN=dev.appsite.com build you 0.1 dev`


## Prerequisites

.Essix apps are built in the Go language.
.The canonical way of deploying Essix apps relies on Docker.
.The `essix` command runs on bash.
.Essix data is stored in RethinkDB, which the `essix` command completely manages
inside Docker, avoiding any additional dependencies.
