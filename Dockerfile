FROM scratch
MAINTAINER eelco@servicelab.org
COPY tamtam /
ENTRYPOINT ["/tamtam"]
CMD ["help"]
