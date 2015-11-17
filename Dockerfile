FROM scratch
MAINTAINER YI-HUNG JEN <yihungjen@gmail.com>

COPY ca-certificates.crt /etc/ssl/certs/
COPY docker-ambassador /
ENTRYPOINT ["/docker-ambassador"]
CMD ["--help"]
