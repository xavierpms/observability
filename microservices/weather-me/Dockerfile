FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -C "server" -ldflags="-w -s" -o weather-me .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server/weather-me .
USER 65532:65532
CMD ["./weather-me"]
