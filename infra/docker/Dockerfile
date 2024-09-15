FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o api ./cmd/api

###

FROM alpine:3.20.0

LABEL maintainer="Sebastian Salata R-T <sa.salatart@gmail.com>"

WORKDIR /root/

COPY --from=builder /app/api .
COPY --from=builder /app/config/config.yaml ./config/config.yaml
COPY --from=builder /app/db/migrations/*.sql ./db/migrations/
COPY --from=builder /app/http/oapi/schema.yaml ./http/oapi/schema.yaml

EXPOSE 8080

CMD ["./api", "-generate"]
