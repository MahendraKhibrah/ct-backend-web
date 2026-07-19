# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/ct-backend .

FROM alpine:3.21 AS runtime
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app
COPY --from=builder --chown=app:app /out/ct-backend ./ct-backend

ENV APP_ENV=production \
    GIN_MODE=release \
    GOLANG_PORT=8888

EXPOSE 8888
USER app

ENTRYPOINT ["/app/ct-backend"]
