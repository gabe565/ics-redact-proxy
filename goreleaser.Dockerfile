FROM gcr.io/distroless/static:nonroot
LABEL org.opencontainers.image.source="https://github.com/gabe565/ics-redact-proxy"
WORKDIR /
COPY ics-redact-proxy /
ENTRYPOINT ["/ics-redact-proxy"]
