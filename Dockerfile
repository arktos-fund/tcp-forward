FROM ubuntu:latest AS rootless

RUN groupadd rootless && \
    useradd -u 10001 noroot && \
    usermod -a -G rootless noroot
RUN apt-get update
RUN apt-get install -y ca-certificates

FROM scratch

COPY --from=rootless /etc/passwd /etc/passwd
COPY --from=rootless /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chown=noroot go /go
USER noroot

VOLUME [ "/config/default.yaml" ]

ENTRYPOINT ["/go"]