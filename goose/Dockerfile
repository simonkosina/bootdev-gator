FROM gomicro/goose:3.25.0

RUN apk update && \
    apk add --no-cache bash postgresql-client

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
