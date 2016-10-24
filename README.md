# Essix
Package essix runs an essential simple secure stable scalable stateless server.

The `essix` [command](#essix-command) manages Essix apps, their TLS (https) certificates, their
backend databases (RethinkDB), and their infrastructure (Docker Swarm Mode).

Follow the [Quickstart](#quickstart) to get your first app running within minutes, on a swarm near you.

With `$essix init`, a new Essix app package is inititialised, much like
https://github.com/wscherphof/essix/tree/master/app, where you would add your
own functionality. It includes the Profile example of how things can be done.

[![GoDoc](https://godoc.org/github.com/wscherphof/essix?status.svg)](https://godoc.org/github.com/wscherphof/essix)

## Features

Essix basically provides what you _need_ for running a reliable web application,
cutting out the cruft, and aiming for a smooth and transparent developer
experience whith the lowest possible rate of surprises.


- [HTTP/2](https://en.wikipedia.org/wiki/HTTP/2)
- Encrypted communication (https)
- One-line commands for generating trusted server certificates using [LetsEncrypt](https://letsencrypt.org/)
- Redundant distributed backend database using [RethinkDB](https://www.rethinkdb.com/)
- [Stateless](https://en.wikipedia.org/wiki/Stateless_protocol) distributed application server
- One-line commands for deploying on [Docker Swarm Mode](https://docs.docker.com/engine/swarm/) computing clusters
- One-line commands for scaling by ading/removing application server or database instances
- One-line commands for scaling by ading/removing swarm computing nodes
- Development environment identical to production environment
- Transparent Swarm environments ranging from a one-node laptop to multi-continent clouds


- Secure user sessions (log in, log out)
- Sign up with email address verification
- Secure password reset
- Secure email address changing process
- Built-in email sending capability
- Built-in i18n multi-language "messages" system


- Automatic [CSRF](https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)) protection using form tokens
- Built-in request route rate limiting (obsoleting the need for [captchas](https://www.owasp.org/index.php/Testing_for_Captcha_(OWASP-AT-012)#WARNING:_CAPTCHA_protection_is_an_ineffective_security_mechanism_and_should_be_perceived_as_a_.22rate_limiting.22_protection_only.21))


- Straightforward business data entity modeling & manipulation
- [Jade](http://jadelang.net/)-like HTML templating, using [Ace](https://github.com/yosssi/ace)
- Built-in error handling
- [Post/Redirect/Get](https://en.wikipedia.org/wiki/Post/Redirect/Get) made simple


## Quickstart

1. Verify the [Prerequisites](#prerequisites).
1. Install Essix: `$ go get -u github.com/wscherphof/essix`
1. Initialise a new Essix app: `$ essix init github.com/you/yourapp`
1. Create a self-signed [TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) certificate: `$ cd $GOPATH/github.com/you/yourapp`, then: `$ essix cert dev.appsite.com`
1. Create a local one-node Docker swarm: `$ essix nodes -H dev.appsite.com -m 1 create dev`
1. Install RethinkDB on the swarm: `$ essix r create dev`
1. Run your app on the swarm: `$ cd $GOPATH/github.com/you/yourapp`, then: `$ essix -e DOMAIN=dev.appsite.com build you 0.1 dev`
1. Point your browser to https://dev.appsite.com/. It'll complain about not trusting your self-signed certificate, but you can instruct it to accept it anyway. The `essix` command can generate officially trusted certificates as well.
1. Put your app `$GOPATH/github.com/you/yourapp` under [version control](https://guides.github.com/introduction/getting-your-project-on-github) & get creative.

Run `$ essix help` for some more elaborate [usage examples]((#essix-command)).

## Prerequisites

1. Essix apps are built in the [Go language](https://golang.org/doc/install).
1. Do create a Go working directory & set the `$GOPATH` environment variable to point to it.
1. The `go get` command relies on a [version control system]().
1. The `essix` command runs on bash, which should be present on Mac or Linux.
For Windows, [Git Bash](https://git-for-windows.github.io/) should work.
1. Deploying Essix apps with the `essix` command, relies on [Docker](https://www.docker.com/products/docker).
As well as on Docker Machine & VirtualBox, which are normally included in a Docker installation.

Essix data is stored in [RethinkDB](https://www.rethinkdb.com/), which the
`essix` command completely manages within Docker, so not a prerequisite.

After installing Essix with `$ go get -u github.com/wscherphof/essix`, the
`essix` command relies on the `$GOPATH/github.com/wscherphof/essix` directory;
leave it untouched.

## Essix command
```
Usage:

essix init PACKAGE
  Initialise a new essix app in the given new directory under $GOPATH.

essix cert DOMAIN [EMAIL]
  Generate a TLS certificate for the given domain.
  Certificate gets saved in ./resources/certificates
  Without EMAIL,
    a self-signed certificate is produced.
  With EMAIL,
    a trusted certificate is produced through LetsEncrypt.
    The current LetsEncrypt approach relies on a DNS configuration on DigitalOcean,
    and requires the DIGITALOCEAN_ACCESS_TOKEN environment variable.

essix nodes [OPTIONS] COMMAND SWARM
  Create a Docker Swarm Mode swarm, and manage its nodes.
  Run 'essix nodes help' for more information.

essix r [OPTIONS] [COMMAND] SWARM
  Create a RethinkDB cluster on a swarm, and/or start its web admin.
  Run 'essix r help' for more information.

essix [OPTIONS] build REPO TAG [SWARM]
  Format & compile go sources in the current directory, and build a Docker
  image named REPO/APP:TAG
  APP is the current directory's name, is also the service name.
  Without SWARM,
    the OPTIONS are ignored, and the image is built locally,
    then pushed to the repository. Default repository is Docker Hub.
  With SWARM,
    the image is built remotely on each of the swarm's nodes,
    and the service is run there, with the given OPTIONS.

essix [OPTIONS] run REPO TAG SWARM
  Run a service from an image on a swarm.
  Options:
    -e key=value ...  environment variables
    -r replicas       number of replicas to run (default=1)

essix help
  Display this message.


You'll want to have these baseline tools ready:
  - Bash
  - Git
  - Go
  - Docker
  - Docker Machine
  - VirtualBox
  - An account with Docker Hub
  - An account with DigitalOcean


Examples:

  $ essix init github.com/essix/newapp
      Initialises a base structure for an Essix app in /Users/wsf/go/github.com/essix/newapp.

  $ essix cert dev.appsite.com
      Generates a self-signed TLS certificate for the given domain.

  $ export DIGITALOCEAN_ACCESS_TOKEN="94dt7972b863497630s73012n10237xr1273trz92t1"
  $ essix cert www.appsite.com essix@appsite.com
      Generates a trusted TLS certificate for the given domain.

  $ essix nodes -m 1 -w 2 -H dev.appsite.com create dev
      Creates swarm dev on VirtualBox, with one manager node, and 2 worker
      nodes. Adds hostname dev.appsite.com to /etc/hosts, resolving to the
      manager node's ip address.

  $ essix nodes -m 1 -d digitalocean -F create www
      Creates a one-node swarm www on DigitalOcean, and enables a firewall on it.
  $ essix nodes -w 1 -d digitalocean -F create www
      Adds a worker node (with firewall) to swarm www on DigitalOcean.

  $ essix r create dev
      Creates a RethinkDB cluster on swarm dev, and opens the cluster's
      administrator web page.

  $ essix r dev
      Opens the dev swarm RethinkDB cluster's administrator web page.

  $ essix build essix 0.2
      Locally builds the essix/APP:0.2 image, and pushes it to the repository.
  $ essix run -e DOMAIN=www.appsite.com -r 6 essix 0.2 www
      Starts 6 replicas of the service on swarm www, using image essix/APP:0.2,
      which is downloaded from the repository, if not found locally.

  $ essix -e DOMAIN=dev.appsite.com build essix 0.3 dev
      Builds image essix/APP:0.3 on swarm dev's nodes, and runs the service
      on dev, with the given DOMAIN environment variable set.
```
