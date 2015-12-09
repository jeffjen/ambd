FROM scratch
MAINTAINER YI-HUNG JEN <yihungjen@gmail.com>

COPY ca-certificates.crt /etc/ssl/certs/
COPY ambd /
COPY ambctl/ambctl /
ENTRYPOINT ["/ambd"]
CMD ["--help"]

ENV VERSION latest
ENV BUILD golang-1.5.1

ENV NODE_NAME ""
ENV NODE_AVAIL_ZONE ""
ENV NODE_REGION ""
ENV NODE_PUBLIC_IPV4 ""
ENV NODE_PRIVATE_IPV4 ""
