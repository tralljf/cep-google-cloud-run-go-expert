FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:3.22

WORKDIR /app

RUN adduser -D -H appuser
COPY --from=builder /server /server

ENV PORT=8080
EXPOSE 8080

USER appuser
ENTRYPOINT ["/server"]
