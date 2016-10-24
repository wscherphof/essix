# Essix
Package essix runs an essential simple secure stable scalable stateless server.

The `essix` command manages Essix apps, their TLS (https) certificates, their backend databases (RethinkDB), and their infrastructure (Docker Swarm Mode).

Follow the Quickstart to get your first app running on a swarm within minutes.

## Quickstart

1. Install Essix: `$ go get -u github.com/wscherphof/essix`
1. Initialise a new Essix app: `$ essix init github.com/you/yourapp`
1. Create a self-signed [TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) certificate: `$ cd $GOPATH/github.com/you/yourapp`, then: `$ essix cert dev.appsite.com`
1. Create a local one-node Docker swarm: `$ essix nodes -H dev.appsite.com -m 1 create dev`
1. Install RethinkDB on the swarm: `$ essix r create dev`
1. Run your app on the swarm: `$ cd $GOPATH/github.com/you/yourapp`, then: `$ essix -e DOMAIN=dev.appsite.com build you 0.1 dev`
1. Point your browser to https://dev.appsite.com/. It'll complain about not trusting your self-signed certificate, but you can instruct it to accept it anyway. The `essix` command can generate officially trusted certificates as well.
1. Put your app `$GOPATH/github.com/you/yourapp` under [version control](https://guides.github.com/introduction/getting-your-project-on-github) & get creative.

Run `$ essix help` for some more elaborate usage examples.

## Prerequisites

1. Essix apps are built in the [Go language](https://golang.org/doc/install).
1. Do create a Go working directory & set the `$GOPATH` environment variable to point to it.
1. The `go get` command relies on a [version control system]().
1. The `essix` command runs on bash, which should be present on Mac or Linux.
For Windows, [Git Bash](https://git-for-windows.github.io/) should work.
1. Deploying Essix apps with the `essix` command, relies on [Docker](https://www.docker.com/products/docker).
As well as on Docker Machine & VirtualBox, which are normally included in a Docker installation.

Essix data is stored in [RethinkDB](https://www.rethinkdb.com/), which the `essix` command completely manages
within Docker, so not a prerequisite.

After installing Essix with `$ go get -u github.com/wscherphof/essix`, the `essix` command relies on the `$GOPATH/github.com/wscherphof/essix` directory; leave it untouched.
