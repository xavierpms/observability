FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -C "cmd/server" -ldflags="-w -s" -o weather-by-city .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cmd/server/weather-by-city .
USER 65532:65532
CMD ["./weather-by-city"]
