# build stage
FROM golang:1.10 AS build-env
MAINTAINER eelco@servicelab.org
ARG git_commit_sha
RUN curl https://glide.sh/get | sh
COPY . /go/src/github.com/servicelab/tamtam
WORKDIR /go/src/github.com/servicelab/tamtam
RUN ./dist.sh linux "$git_commit_sha"

# final stage
FROM scratch
MAINTAINER eelco@servicelab.org
WORKDIR /app
COPY --from=build-env /go/src/github.com/servicelab/tamtam/tamtam /app/
ENTRYPOINT ./tamtam
