FROM alpine:latest

# Install time zone package
RUN apk add --update-cache tzdata
# Remove tempfile and apk command.
RUN rm -rf /tmp/* /var/cache/apk/* /sbin/apk

# Copy binary
WORKDIR /blockhw
COPY build/indexer /blockhw/indexer
COPY ./config /blockhw/config
COPY ./docker/wait-for-it.sh /blockhw/wait-for-it.sh

ENTRYPOINT ["/blockhw/indexer"]
CMD ["./wait-for-it.sh", "localhost:5431", "--", "/blockhw/indexer"]