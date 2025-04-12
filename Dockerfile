# Builder
FROM golang:1.24-alpine AS builder

ARG VERSION=dev
ARG BUILD_DATE
ARG COMMIT_REF

LABEL\
	maintainer="Causality <mail@causality.africa" \
	org.opencontainers.image.authors="Causality <mail@causality.africa>" \
	org.opencontainers.image.created=$BUILD_DATE \
	org.opencontainers.image.description="Golang backend for Causality Africa" \
	org.opencontainers.image.documentation="https://causality.africa" \
	org.opencontainers.image.licenses="MIT" \
	org.opencontainers.image.revision=$COMMIT_REF \
	org.opencontainers.image.source="https://github.com/causality-africa/core" \
	org.opencontainers.image.title="Causality Core" \
	org.opencontainers.image.url="https://github.com/causality-africa/core/pkgs/container/core" \
	org.opencontainers.image.vendor="Causality" \
	org.opencontainers.image.version=$VERSION

ENV CGO_ENABLED=1

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/jackc/tern/v2@latest

COPY . .

RUN go build \
    -ldflags="-w -s -X main.Version=${VERSION}" \
    -o core ./cmd/core

# core image
FROM alpine:latest

COPY --from=builder /app/core /core
COPY --from=builder /app/migrations /migrations
COPY --from=builder /go/bin/tern /tern
COPY --from=builder /app/tern.conf /tern.conf

EXPOSE 8080

CMD ["/core"]
