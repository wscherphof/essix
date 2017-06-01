# jmeter

Manage Apache JMeter (distributed) load tests.

<pre>

Usage:

jmeter [OPTONS] run JMX
  Run the given test in non-gui mode, locally or remote, generating
  a dashboard report.
  Prerequitites:
    - JMeter 3.2 installed: http://jmeter.apache.org/
    - Environment variable JMETER_HOME set (default: ~/apache-jmeter-3.2)
  JMX is the path to a JMeter test plan definition .jmx file.
  Options:
    -i ip_address    IP address to use for the master (client)
                     Environment variable IP works as well.
                     If unset, Java picks one.
    -R remote_hosts  List of slave hosts to use for distributed testing.
                     Environment variable REMOTE_HOSTS works as well.
                     If unset, test is run locally.
    -d dir           Base directory for test results.
                     Default: ./jmeter-test

jmeter server ACTION MACHINE
  Manage remote JMeter slaves (servers)
  Prerequitites:
    - Docker & Docker Machine installed: https://docs.docker.com/machine/
  ACTION is either start, stop, or restart.
  MACHINE is the targeted docker-machine.

jmeter perfmon ACTION MACHINE
  Manage the PerfMon Server Agent on application servers.
  Prerequitites:
    - Docker & Docker Machine installed: https://docs.docker.com/machine/
  ACTION is either start, stop, or restart.
  MACHINE is the targeted docker-machine.

jmeter help
  Display this message.

</pre>


