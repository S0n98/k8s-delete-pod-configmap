FROM golang:1.21.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /out/main ./

FROM alpine
WORKDIR /out/mainout/main
COPY --from=builder /out/main /out/main
ENTRYPOINT ["/out/main"]