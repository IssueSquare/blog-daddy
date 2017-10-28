# final stage
FROM centurylink/ca-certs
COPY bin/app /app
COPY app/config.yaml /
ENTRYPOINT ["/app"]
