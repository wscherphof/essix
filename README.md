# Essix
Package essix runs an essential simple secure stable scalable stateless server.

The `essix` command manages Essix apps, their TLS (https) certificates, their
backend databases (RethinkDB), and their infrastructure (Docker Swarm Mode).

Follow the [Quickstart](#quickstart) to get your first app running on a swarm within minutes.

With `$essix init`, a new Essix app package is inititialised, much like
https://github.com/wscherphof/essix/tree/master/app, where you would add your
own functionality. It includes the Profile example of how things can be done.

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
- Transparent Swarm environment ranging from a one-node laptop to multi-continent clouds

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

Run `$ essix help` for some more elaborate usage examples.

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
