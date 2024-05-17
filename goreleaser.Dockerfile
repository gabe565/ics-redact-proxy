FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY ics-availability-server /
ENTRYPOINT ["/ics-availability-server"]
