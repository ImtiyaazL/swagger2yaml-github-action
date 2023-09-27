FROM golang:1.20 AS builder

WORKDIR /build

COPY . .

RUN go build -o bin/swagger2yaml ./v2

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV ACCOUNT=""
ENV HOST=""
ENV INPUT_FILE=""
ENV OUTPUT_FILE=""
ENV REGION=""
ENV VPC=""

ENTRYPOINT ["/entrypoint.sh"]
