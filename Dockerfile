FROM alpine

ADD bin/linux/* /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/fluxd"]
