FROM scratch
MAINTAINER eelco@servicelab.org
ARG BIN=tamtam
WORKDIR /app
COPY $BIN /app/runner
ENTRYPOINT ["/app/runner"]
CMD ["help"]
