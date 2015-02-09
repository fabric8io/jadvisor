# jAdvisor [![Circle CI](https://circleci.com/gh/fabric8io/jadvisor.svg?style=svg)](https://circleci.com/gh/fabric8io/jadvisor)

jAdvisor is here to help you gather metrics from your [Jolokia](http://jolokia.org/) enabled
[Docker](https://docker.com/) containers running in your Kubernetes environment.

## Running

```
Usage of ./stage/jadvisor:
  -alsologtostderr=false: log to standard error as well as files
  -kubernetes_insecure=false: Trust Kubernetes master certificate (if using https)
  -kubernetes_master="https://localhost:8443": Kubernetes master address
  -log_backtrace_at=:0: when logging hits line file:N, emit a stack trace
  -log_dir="": If non-empty, write log files in this directory
  -logtostderr=false: log to standard error instead of files
  -poll_duration=10s: Polling duration
  -sink="memory": Backend storage. Options are [memory | influxdb]
  -sink_influxdb_buffer_duration=10s: Time duration for which stats should be buffered in influxdb sink before being written as a single transaction
  -sink_influxdb_host="localhost:8086": InfluxDB host:port
  -sink_influxdb_name="k8s": Influxdb database name
  -sink_influxdb_password="root": InfluxDB password
  -sink_influxdb_username="root": InfluxDB username
  -sink_memory_ttl=1h0m0s: Time duration for which stats should be cached if the memory sink is used
  -stderrthreshold=0: logs at or above this threshold go to stderr
  -v=0: log level for V logs
  -vmodule=: comma-separated list of pattern=N settings for file-filtered logging
```

## Building

First, ensure your `GOPATH` environment variable is set up properly. You can then get this
source repo with:

```
go get -d github.com/fabric8io/jadvisor
```

This will clone the jadvisor into `$GOPATH/src/github.com/fabric8io/jadvisor`.

You will also need to install [Godep](https://github.com/tools/godep) by running:

```
go get github.com/tools/godep
```

Ensure that `$GOPATH`/bin is in your path in order to use `godep`.

Now change into the jadvisor source dir:

```
cd $GOPATH/src/github.com/fabric8io/jadvisor/
```

To build the software run:

```
make
```

This will build the binary into `stage` directory & also build a Docker image `jadvisor:latest` for you to try out.
