#FROM alpine
FROM ubuntu:xenial

COPY bin/huepro /bin/huepro

ENTRYPOINT ["/bin/huepro"]
