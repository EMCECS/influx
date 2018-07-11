FROM alpine

RUN apk add libc6-compat
ADD bin/linux/fluxd /usr/local/bin/fluxd

ENTRYPOINT /usr/local/bin/flux