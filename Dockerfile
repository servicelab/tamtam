# build stage
FROM golang:1.10 AS build-env
MAINTAINER eelco@servicelab.org
ARG git_commit_sha
COPY . /go/src/github.com/servicelab/tamtam
WORKDIR /go/src/github.com/servicelab/tamtam
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh| sh \
#    && dep ensure \
    && make docker

# final stage
FROM scratch
MAINTAINER eelco@servicelab.org
WORKDIR /app
COPY --from=build-env /go/src/github.com/servicelab/tamtam/dist/tamtam /app/
CMD ["./tamtam"]
