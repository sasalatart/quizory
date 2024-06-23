FROM golang:1.22-alpine AS builder

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
COPY --from=builder /app/db/migrations/*.sql ./db/migrations/
COPY --from=builder /app/http/oapi/schema.yml ./http/oapi/schema.yml

ENV DB_URL=$DB_URL
ENV OPENAI_API_KEY=$OPENAI_API_KEY
ENV JWT_SECRET=$JWT_SECRET

EXPOSE 8080

CMD ["./api", "-generate"]
