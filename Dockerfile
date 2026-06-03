FROM golang:1.26 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mediadex ./main.go

FROM scratch
COPY --from=builder /app/mediadex /mediadex

EXPOSE 8080
ENTRYPOINT ["/mediadex"]
