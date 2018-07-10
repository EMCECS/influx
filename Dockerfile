FROM alpine

ADD bin/linux/fluxd /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/fluxd"]
