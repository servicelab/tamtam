FROM alpine
MAINTAINER eelco@servicelab.org
ARG BIN=
WORKDIR /app
COPY $BIN /app/runner
ENTRYPOINT ["/app/runner"]
CMD ["help"]
