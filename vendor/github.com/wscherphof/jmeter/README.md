# jmeter

Manage Apache JMeter (distributed) load tests with Docker.

In other words: JMeter [MTA](https://blog.docker.com/2017/04/modernizing-traditional-apps-with-docker/)'d :-)

```
Prerequitites:
  - Docker & Docker Machine installed: https://docs.docker.com/machine/
  - A test plan definition .jmx file created with JMeter 3.2: http://jmeter.apache.org/

Usage:

jmeter [OPTIONS] run JMX MACHINE [REMOTE_MACHINE...]
  Run the given test in non-gui mode, locally or remote, generating
  a dashboard report.
  JMX                 Path to the .jmx file.
  MACHINE             The docker-machine that should run the test.
  REMOTE_MACHINE...   Any number of names of docker-machines with remote
                      slave servers to use for distributed testing.
                      If unset, the test is run locally.
  Options:
    -p key=value ...  Properties; set both locally, on the controller
                      client, as globally, on the slave servers.

jmeter server ACTION MACHINE[...]
  Manage remote JMeter slave servers.
  ACTION   Either start, stop, or restart.
  MACHINE  One or more targeted docker-machines.

jmeter perfmon ACTION MACHINE[...]
  Manage the PerfMon Server Agent on application servers.
  ACTION   Either start, stop, or restart.
  MACHINE  One or more targeted docker-machines.

jmeter help
  Display this message.
```

