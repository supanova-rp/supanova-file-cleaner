# =========================
# Stage 1: Build the binary
# =========================
FROM --platform=$BUILDPLATFORM golang:1.25 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN make dep

COPY . .

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make build

# =========================
# STAGE 2: runtime
# =========================
FROM alpine:3.22

RUN addgroup -g 1000 -S appgroup && adduser -u 1000 -S appuser -G appgroup

RUN mkdir /app && chown appuser:appgroup /app
WORKDIR /app

COPY --from=builder /app/supanova-file-cleaner /app/supanova-file-cleaner

USER appuser

ENTRYPOINT ["/app/supanova-file-cleaner"]
