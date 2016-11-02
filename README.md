# Essix
Essix runs an essential simple secure stable scalable stateless server.

`$ go get -u github.com/wscherphof/essix`

The `essix` [command](#essix-command) manages Essix apps, their server
certificates (TLS/HTTPS), their backend databases (RethinkDB), and their
infrastructure (Docker Swarm Mode).

Follow the [Quickstart](#quickstart) to get your first app running within
minutes, on a swarm near you.

With `$ essix init`, a new Essix app package is inititialised, much like
https://github.com/wscherphof/essix/tree/master/app, where you would add your
own functionality. It includes the Profile example of how things can be done.

[![GoDoc](https://godoc.org/github.com/wscherphof/essix?status.svg)](https://godoc.org/github.com/wscherphof/essix)

## Features
Essix basically provides _just what you need_ for running a reliable web
application.

### Dev
Essix cuts out the cruft, and facilitates building directly on the excellent
standards that created the web.

- The server runs [HTTP/2](https://en.wikipedia.org/wiki/HTTP/2) with
[TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) (https).
- `$ essix cert` generates self-signed _TLS certificates_, and trusted
certficates for domains you own, through [LetsEncrypt](https://letsencrypt.org/)
- Templates for [HTML](https://www.w3.org/html/) documents are defined with
[Ace](https://github.com/yosssi/ace), a Go version of
[Jade](http://jadelang.net/).
[Progressively enhance](https://en.wikipedia.org/wiki/Progressive_enhancement)
them with [CSS](https://www.w3.org/Style/CSS/) styles,
[SVG](https://www.w3.org/Graphics/SVG/) graphics, and/or
[JavaScript](https://www.w3.org/standards/webdesign/script) behaviours as you
like. Custom templates may override core templates.
- Every path in `<app>/resources/static/...` is _statically served_ as
`https://<host>/static/...`
- Request _routes_ are declared with their method, URL, and handler function,
e.g. `router.GET("/account/activate", account.ActivateForm)`
- Business objects gain
[CRUD](https://en.wikipedia.org/wiki/Create,_read,_update_and_delete) operations
from the [Entity](https://godoc.org/github.com/wscherphof/entity) base type,
which manages their storage in a [RethinkDB](https://www.rethinkdb.com/)
cluster.
- Server and user errors are communicated through a customisable
[error](https://godoc.org/github.com/wscherphof/essix/template#Error)
[template](https://github.com/wscherphof/essix/blob/master/resources/templates/template/Error.ace).
- The [Post/Redirect/Get](https://en.wikipedia.org/wiki/Post/Redirect/Get)
pattern is a
[first class citizen](https://godoc.org/github.com/wscherphof/essix/template#PRG).
- HTML _email_ is sent using the same [Ace](https://github.com/yosssi/ace)
[templates](https://godoc.org/github.com/wscherphof/essix/template#Email). Failed
emails are queued automatically to send later.
- Multi-language labels and text are managed through the simple definition of
[messages](https://godoc.org/github.com/wscherphof/msg) with keys and
accompanying translations. Custom messages may override core messages.

### Ops
Essix creates computing environments from scratch in a snap, scales
transparently from a local laptop to a multi-continent cloud, and only knows how
to run in fault tolerant mode.

- `$ essix nodes` creates and manages
[Docker Swarm Mode](https://docs.docker.com/engine/swarm/) swarms, either
locally or in the cloud. See the
[blogpost](https://wscherphof.wordpress.com/2016/09/13/rethink-swarm-mode/)
for the ins and outs.
- `$ essix r` installs a [RethinkDB](https://www.rethinkdb.com/) cluster on a
swarm's nodes.
- `$ essix build` compiles an Essix app's sources, and builds a
[Docker](https://www.docker.com/) image for it.
- `$ essix run` creates a service on the swarm that runs any number of
_replicas_ of the app's image.
- The app server is
[stateless](http://whatisrest.com/rest_constraints/stateless_profile)
(resource data is kept in the database cluster, and user session data is kept
client-side in a cookie), meaning each replica is the same as any of the others,
every request can be handled by any of the replicas, and if one fails, the
others continue to serve.

### Security
All communication between client and server is encrypted through
[TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) (https). User
session tokens are stored as a
[secure cookie](http://www.gorillatoolkit.org/pkg/securecookie)

- HTTP PUT, POST, PATCH, and DELETE requests are protected from
[CSRF](https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF))
attacks automatically, using encrypted
[form tokens](https://godoc.org/github.com/wscherphof/secure#SecureRouter).
- On sign up, the user's _email address_ is verified before the new account is
activated. User _passwords_ are never stored; on sign in, the given password is
verified through an encrypted hash value in the database. The processes for
resetting the password, changing the email address, or suspending an account,
include an _email verification_ step.
- Specific request routes can be declaratively shielded from
[unauthorised access](https://godoc.org/github.com/wscherphof/secure#Handle),
e.g. `router.PUT("/account/email", secure.Handle(account.ChangeEmail))`
- Specific request routes can be declaratively
[rate limited](https://godoc.org/github.com/wscherphof/essix/ratelimit#Handle)
(obsoleting the need for
[captchas](https://www.owasp.org/index.php/Testing_for_Captcha_(OWASP-AT-012)#WARNING:_CAPTCHA_protection_is_an_ineffective_security_mechanism_and_should_be_perceived_as_a_.22rate_limiting.22_protection_only.21))
e.g. `router.PUT("/session", ratelimit.Handle(account.LogIn))`
- A _firewall_ is included for cloud nodes, opening only ports 80, and 443 for
the app (80 redirects to 443), 2376, 2377, 7946, and 4789 for Docker Swarm Mode,
and 22 for `ssh`. An ssh tunnel provides access to the RethinkDB admin site.

## Quickstart

1. Verify the [Prerequisites](#prerequisites).
1. Install Essix: `$ go get -u github.com/wscherphof/essix`
1. Initialise a new Essix app: `$ essix init github.com/you/yourapp`
1. Create a self-signed
[TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) certificate:
`$ cd $GOPATH/src/github.com/you/yourapp`, then: `$ essix cert dev.appsite.com`
1. Create a local one-node Docker swarm: `$ essix nodes -H dev.appsite.com -m 1
create dev`
1. Install RethinkDB on the swarm: `$ essix r create dev`
1. Run your app on the swarm: `$ cd $GOPATH/src/github.com/you/yourapp`, then:
`$ essix -e DOMAIN=dev.appsite.com build you 0.1 dev`
1. Point your browser to https://dev.appsite.com/. It'll complain about not
trusting your self-signed certificate, but you can instruct it to accept it
anyway. `$ essix cert` can generate officially trusted certificates as well.
1. [Configure](#email-configuration) the server's email account details.
1. Put your app (`$GOPATH/src/github.com/you/yourapp`) under
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
`essix` command relies on the `$GOPATH/src/github.com/wscherphof/essix` directory;
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

  $ export DIGITALOCEAN_ACCESS_TOKEN="94dt7972b863497630s73012n10237xr1273trz92t1"
  $ essix nodes -m 1 -d digitalocean -F create www
      Creates one-node swarm www on DigitalOcean, with a firewall enabled.
  $ export DIGITALOCEAN_REGION="ams3"
  $ essix nodes -w 1 -d digitalocean -F create www
      Adds an Amsterdam based worker node to swarm www.
  $ export DIGITALOCEAN_REGION="sgp1"
  $ essix nodes -w 1 -d digitalocean -F create www
      Adds a Singapore based worker node to swarm www.

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

## Details

### Email configuration
The [email](https://godoc.org/github.com/wscherphof/essix/email) function needs
email server & account details configured in the database. `$ essix r dev`
brings up the administrator site of the RethinkDB cluster (assuming you target
the `dev` swarm). Navigate to the Data Explorer, and paste this command:
```
r.db('essix').table('config').get('email').update({
  EmailAddress: 'essix@gmail.com',
  PWD: 'secret',
  PortNumber: '587',
  SmtpServer: 'smtp.gmail.com'
}
```
Replace the values with what's appropriate for your account, then run the
command. To get the app server to read the new config, restart it:
```
$ docker-machine ssh dev-manager-1 docker service scale myapp=0
myapp scaled to 0
$ docker-machine ssh dev-manager-1 docker service scale myapp=1
myapp scaled to 1
```

Note that should you choose to use a Gmail account, you need to turn on 'Allow
Less Secure Apps to Access Account' in the
[account settings](https://myaccount.google.com/u/1/security)
