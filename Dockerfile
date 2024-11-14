FROM golang:1.23.2-alpine3.19

RUN apk add --no-cache postgresql-client

RUN mkdir -p /home/app

COPY . /home/app

WORKDIR /home/app

RUN go mod tidy
RUN go build -o ideyar ./cmd/ideyar/*

COPY scripts/create-tables.sh /home/app/scripts/create-tables.sh

RUN chmod +x /home/app/scripts/create-tables.sh

RUN echo '#!/bin/sh' > /home/app/entrypoint.sh && \
    echo '/home/app/scripts/create-tables.sh' >> /home/app/entrypoint.sh && \
    echo 'exec ./ideyar "$@"' >> /home/app/entrypoint.sh

RUN chmod +x /home/app/entrypoint.sh

ENTRYPOINT ["sh", "/home/app/entrypoint.sh"]
