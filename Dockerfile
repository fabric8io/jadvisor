FROM flynn/busybox
MAINTAINER Jimmi Dyson <jimmidyson@gmail.com>

ADD ./stage/jadvisor /bin/jadvisor

ENTRYPOINT ["/bin/jadvisor"]
CMD []
