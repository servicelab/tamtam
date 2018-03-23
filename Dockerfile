# build stage
FROM golang:1.9 AS build-env
ARG git_commit_sha
RUN curl https://glide.sh/get | sh
ADD . /go/src/github.com/eelcocramer/tamtam
WORKDIR /go/src/github.com/eelcocramer/tamtam
RUN ./dist.sh linux $git_commit_sha

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/eelcocramer/tamtam/tamtam /app/
ENTRYPOINT ./tamtam

