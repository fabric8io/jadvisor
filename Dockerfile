FROM flynn/busybox
MAINTAINER Jimmi Dyson <jimmidyson@gmail.com>

ADD ./build/jadvisor /bin/jadvisor

ENTRYPOINT ["/bin/jadvisor"]
CMD []
