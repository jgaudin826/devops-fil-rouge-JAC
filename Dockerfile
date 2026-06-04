FROM golang:1.26.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /mediadex ./main.go

FROM alpine:3.20 AS runner
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN apk add --no-cache ca-certificates wget

WORKDIR /app
COPY --from=builder /mediadex /mediadex
RUN chown appuser:appgroup /mediadex
RUN chmod +x /mediadex

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8080/ || exit 1

ENTRYPOINT ["/mediadex"]
