# build stage
FROM golang:1.10 AS builder
MAINTAINER eelco@servicelab.org
ARG GOOS=linux
ARG GOARCH=amd64
ARG PACKAGE=
ARG VERSION=
ARG TIME=
ARG HASH=
ENV LDFLAGS="-s -w -X $PACKAGE/cmd.Version=$VERSION -X $PACKAGE/cmd.BuildTime=$TIME -X $PACKAGE/cmd.GitHash=$HASH"
COPY . /go/src/$PACKAGE
WORKDIR /go/src/$PACKAGE
RUN go test ./...
RUN echo "##### Building for $GOOS on $GOARCH" && \
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH \
    go build -a --installsuffix cgo -ldflags="$LDFLAGS" -o /go/bin/runner

# final stage
FROM scratch
MAINTAINER eelco@servicelab.org
WORKDIR /app
COPY --from=builder /go/bin/runner /app/
ENTRYPOINT ["/app/runner"]
CMD ["help"]
