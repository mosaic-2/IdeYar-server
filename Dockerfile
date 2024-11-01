FROM golang:1.22.2-alpine3.19

RUN mkdir -p /home/app

COPY . /home/app

WORKDIR /home/app

RUN go mod tidy
RUN go build -o main ./cmd/server/main.go

COPY scripts/create-tables.sh /home/app/scripts/create-tables.sh

RUN chmod +x /home/app/scripts/create-tables.sh

RUN echo '#!/bin/sh\n' > /home/app/entrypoint.sh && \
    echo '/home/app/scripts/create-tables.sh\n' >> /home/app/entrypoint.sh && \
    echo 'exec ./main "$@"\n' >> /home/app/entrypoint.sh

RUN chmod +x /home/app/entrypoint.sh

ENTRYPOINT ["/home/app/entrypoint.sh"]