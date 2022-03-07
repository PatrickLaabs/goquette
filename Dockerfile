# Dockerfile
FROM alpine
COPY goquette /usr/bin/goquette
ENTRYPOINT ["/usr/bin/goquette"]