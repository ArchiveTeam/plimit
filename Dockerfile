FROM golang:1.23.5-bookworm
WORKDIR /plimit/
COPY . .
RUN make

FROM debian:bookworm-slim

RUN apt-get update && \
	apt-get -y upgrade && \
	apt-get -y install --no-install-recommends ca-certificates && \
	apt-get clean && \
	rm -rf /var/lib/apt/lists/*

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

COPY --from=0 /plimit/plimit /plimit

ENTRYPOINT ["/tini", "--", "/plimit"]
CMD ["exporter"]

