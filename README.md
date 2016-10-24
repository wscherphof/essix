# Essix
Package essix runs an essential simple secure stable scalable stateless server.

## Quickstart

1. Install Essix: `$ go get -u github.com/wscherphof/essix`
1. Initialise a new Essix app: `$ essix init github.com/you/yourapp`
1. Create a self-signed TLS certificate: `$ cd $GOPATH/github.com/you/yourapp`, then: `$ essix cert dev.appsite.com`
1. Create a local one-node Docker swarm: `$ essix nodes -H dev.appsite.com -m 1 create dev`
1. Install RethinkDB on the swarm: `$ essix r create dev`
1. Run your app on the swarm: `$ cd $GOPATH/github.com/you/yourapp`, then: `$ essix -e DOMAIN=dev.appsite.com build you 0.1 dev`

Run `$ essix help` for some more elaborate usage examples.

## Prerequisites

- Essix apps are built in the [Go language](https://golang.org/doc/install).
- The `essix` command runs on bash, which should be present on Mac or Linux.
For Windows, [Git Bash](https://git-for-windows.github.io/) should work.
- Deploying Essix apps with the `essix` command, relies on [Docker](https://www.docker.com/products/docker).
As well as on Docker Machine & VirtualBox, which are normally included in a Docker installation.

Essix data is stored in [RethinkDB](https://www.rethinkdb.com/), which the `essix` command completely manages
within Docker, so not a prerequisite.

After installing Essix with `$ go get -u github.com/wscherphof/essix`, the `essix` command relies on the `$GOPATH/github.com/wscherphof/essix` directory, so leave it untouched.
