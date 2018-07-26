# build stage
FROM golang:1.10 AS builder
MAINTAINER eelco@servicelab.org
ARG GOOS=linux
ARG GOARCH=amd64
COPY . /go/src/github.com/servicelab/tamtam
WORKDIR /go/src/github.com/servicelab/tamtam
RUN go test ./...
RUN echo "##### Building for $GOOS on $GOARCH" && \
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH \
    go build -a --installsuffix cgo -ldflags="-w -s" -o /go/bin/runner

# final stage
FROM scratch
MAINTAINER eelco@servicelab.org
WORKDIR /app
COPY --from=builder /go/bin/runner /app/
CMD ["/app/runner"]
