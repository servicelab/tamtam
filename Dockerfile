FROM scratch
MAINTAINER eelco@servicelab.org
ARG BIN=
WORKDIR /app
COPY $BIN /app/tamtam
EXPOSE 9999
EXPOSE 6262
ENTRYPOINT ["/app/tamtam"]
CMD ["help"]
