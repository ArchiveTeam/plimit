FROM golang:1.23.5-bookworm
WORKDIR /plimit/
COPY . .
RUN make

FROM debian:buster-slim

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

COPY --from=0 /plimit/plimit /plimit

ENTRYPOINT ["/tini", "--", "/plimit"]
CMD ["exporter"]

