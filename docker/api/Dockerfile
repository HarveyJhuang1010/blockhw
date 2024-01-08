FROM alpine:latest

# Install time zone package
RUN apk add --update-cache tzdata
# Remove tempfile and apk command.
RUN rm -rf /tmp/* /var/cache/apk/* /sbin/apk

# Copy binary
WORKDIR /blockhw
COPY build/api /blockhw/api
COPY ./config /blockhw/config
COPY ./docker/wait-for-it.sh /blockhw/wait-for-it.sh

ENTRYPOINT ["/blockhw/api"]
CMD ["./wait-for-it.sh", "localhost:5431", "--", "/blockhw/api"]