FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY ics-redact-proxy /
ENTRYPOINT ["/ics-redact-proxy"]
