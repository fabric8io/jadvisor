# jAdvisor

jAdvisor is here to help you gather metrics from your [Jolokia](http://jolokia.org/) enabled
[Docker](https://docker.com/) containers running in your Kubernetes environment.

## Building

First, ensure your `GOPATH` environment variable is set up properly. You can then get this
source repo with:

```
go get -d github.com/fabric8io/jadvisor
```

This will clone the jadvisor into `$GOPATH/src/github.com/fabric8io/jadvisor`.

You will also need to install [Godep](https://github.com/tools/godep)
