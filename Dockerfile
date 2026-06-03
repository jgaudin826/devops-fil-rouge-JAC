FROM golang:1.26.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /mediadex ./main.go

FROM scratch
COPY --from=builder /mediadex /mediadex

EXPOSE 8080
ENTRYPOINT ["/mediadex"]
