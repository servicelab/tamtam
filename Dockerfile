# build stage
FROM golang:1.10 AS builder
MAINTAINER eelco@servicelab.org
ARG git_commit_sha
COPY . /go/src/github.com/servicelab/tamtam
WORKDIR /go/src/github.com/servicelab/tamtam
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/tamtam

# final stage
FROM scratch
MAINTAINER eelco@servicelab.org
WORKDIR /app
COPY --from=builder /go/bin/tamtam /app/
CMD ["/app/tamtam"]
