FROM golang:1.20 AS builder
COPY . /app
WORKDIR /app
RUN go build -tags=jsoniter -o mah-moneh ./cmd/api


FROM debian:bookworm-slim
RUN apt-get update && apt-get install --yes ca-certificates
COPY --from=builder /app/mah-moneh /mah-moneh
CMD ["/mah-moneh"]
