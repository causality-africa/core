# Builder
FROM golang:1.23-alpine AS builder

LABEL\
	maintainer="Causality <mail@causality.africa" \
	org.opencontainers.image.authors="Causality <mail@causality.africa>" \
	org.opencontainers.image.created=$BUILD_DATE \
	org.opencontainers.image.description="Golang backend for Causality Africa" \
	org.opencontainers.image.documentation="https://www.stephenmwangi.com/beacon-dns/" \
	org.opencontainers.image.licenses="MIT" \
	org.opencontainers.image.revision=$COMMIT_REF \
	org.opencontainers.image.source="https://github.com/causality-africa/core" \
	org.opencontainers.image.title="Causality Core" \
	org.opencontainers.image.url="https://github.com/causality-africa/core" \
	org.opencontainers.image.vendor="Causality Africa" \
	org.opencontainers.image.version=$VERSION


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build \
    -ldflags="-w -s -X main.Version=${VERSION}" \
    -o core ./cmd/core

# core image
FROM alpine:latest

COPY --from=builder /app/core /core

CMD ["/core"]