
# The build stage
FROM golang:1.23.3 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/app/*.go

# The run stage
FROM scratch
WORKDIR /app
# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
