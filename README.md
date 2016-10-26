# Essix
Essix runs an essential simple secure stable scalable stateless server.

`$ go get -u github.com/wscherphof/essix`

The `essix` [command](#essix-command) manages Essix apps, their server
certificates (TLS/HTTPS), their backend databases (RethinkDB), and their
infrastructure (Docker Swarm Mode).

Follow the [Quickstart](#quickstart) to get your first app running within
minutes, on a swarm near you.

With `$essix init`, a new Essix app package is inititialised, much like
https://github.com/wscherphof/essix/tree/master/app, where you would add your
own functionality. It includes the Profile example of how things can be done.

[![GoDoc](https://godoc.org/github.com/wscherphof/essix?status.svg)](https://godoc.org/github.com/wscherphof/essix)

## Features
Essix basically provides _just what you need_ for running a reliable web
application.

### Dev
Essix cuts out the cruft, and facilitates building directly on the excellent
standards that created the web.

The server runs [HTTP/2](https://en.wikipedia.org/wiki/HTTP/2) with
[TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) (https).

`$ essix cert` generates self-signed _TLS certificates_, and trusted certficates
for domains you own, through [LetsEncrypt](https://letsencrypt.org/)

Templates for [HTML](https://www.w3.org/html/) documents are defined with
[Ace](https://github.com/yosssi/ace), a Go version of
[Jade](http://jadelang.net/).
[Progressively enhance](https://en.wikipedia.org/wiki/Progressive_enhancement)
them with [CSS](https://www.w3.org/Style/CSS/) styles,
[SVG](https://www.w3.org/Graphics/SVG/) graphics, and/or
[JavaScript](https://www.w3.org/standards/webdesign/script) behaviours as you
like. Custom templates may override core templates.

Business objects gain Create, Read, Update, and Delete operations from the
_Entity base type_, which manages their storage in a
[RethinkDB](https://www.rethinkdb.com/) cluster.

Server errors are communicated through a customisable _error template_.

The [Post/Redirect/Get](https://en.wikipedia.org/wiki/Post/Redirect/Get) pattern
is a first class citizen.

HTML _email_ is sent using using the same [Ace](https://github.com/yosssi/ace)
templates. Failed emails are queued automatically to send later.

Multi-language labels and text are managed through the simple definition of
_messages_ with keys and translations. Custom messages may override core
messages.


### Security
All communication between client and server is encrypted through
[TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) (https).

HTTP PUT, POST, PATCH, and DELETE requests are protected from
[CSRF](https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF))
attacks automatically, using encrypted form tokens.

On sign up, the user's _email address_ is verified before the new account is
activated. User _passwords_ are never stored; on sign in, the given password is
verified through an encrypted hash value in the database. The processes for
resetting the password, changing the email address, or suspending an account,
include an _email verification_ step.

Specific request routes can be declaratively _rate limited_ (obsoleting the need
for [captchas](https://www.owasp.org/index.php/Testing_for_Captcha_(OWASP-AT-012)#WARNING:_CAPTCHA_protection_is_an_ineffective_security_mechanism_and_should_be_perceived_as_a_.22rate_limiting.22_protection_only.21))


### Ops
Essix creates computing environments from scratch in a snap, scales
transparently from a local laptop to a multi-continent cloud, and only knows how
to run in fault tolerant mode.

`$ essix nodes` creates and manages
[Docker Swarm Mode](https://docs.docker.com/engine/swarm/) swarms, either
locally or in the cloud.

`$ essix r` installs a [RethinkDB](https://www.rethinkdb.com/) cluster on a
swarm's nodes.

`$ essix build` compiles an Essix app's sources, and builds a
[Docker](https://www.docker.com/) image for it.

`$ essix run` creates a service on the swarm that runs any number of _replicas_
of the app's image.

The app server is
[stateless](http://whatisrest.com/rest_constraints/stateless_profile)
(resource data is kept in the database cluster, and user session data is kept
client-side in a cookie), meaning each replica is the same as any of the others,
every request can be handled by any of the replicas, and if one fails, the
others continue to serve.


## Quickstart

1. Verify the [Prerequisites](#prerequisites).
1. Install Essix: `$ go get -u github.com/wscherphof/essix`
1. Initialise a new Essix app: `$ essix init github.com/you/yourapp`
1. Create a self-signed
[TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) certificate:
`$ cd $GOPATH/github.com/you/yourapp`, then: `$ essix cert dev.appsite.com`
1. Create a local one-node Docker swarm: `$ essix nodes -H dev.appsite.com -m 1
create dev`
1. Install RethinkDB on the swarm: `$ essix r create dev`
1. Run your app on the swarm: `$ cd $GOPATH/github.com/you/yourapp`, then:
`$ essix -e DOMAIN=dev.appsite.com build you 0.1 dev`
1. Point your browser to https://dev.appsite.com/. It'll complain about not
trusting your self-signed certificate, but you can instruct it to accept it
anyway. `$ essix cert` can generate officially trusted certificates as well.
1. Put your app (`$GOPATH/github.com/you/yourapp`) under
[version control](https://guides.github.com/introduction/getting-your-project-on-github)
& get creative.

Run `$ essix help` for some more elaborate [usage examples]((#essix-command)).

## Prerequisites

1. Essix apps are built in the [Go language](https://golang.org/doc/install).
1. Do create a Go working directory & set the `$GOPATH`
[environment variable](https://golang.org/doc/install#testing) to point to it.
1. The `go get` command relies on a version control system, e.g.
[GitHub](https://github.com/).
1. The `essix` command runs on `bash`, which should be present on Mac or Linux.
For Windows, [Git Bash](https://git-for-windows.github.io/) should work.
1. Deploying apps with the `essix` command, relies on
[Docker](https://www.docker.com/products/docker). Docker Machine & VirtualBox
are normally included in the Docker installation.

Essix data is stored in a [RethinkDB](https://www.rethinkdb.com/) cluster, which
the `essix` command completely manages within Docker, so not a prerequisite.

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


Examples:

  $ essix init github.com/essix/newapp
      Initialises a base structure for an Essix app in $GOPATH/src/github.com/essix/newapp.

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
