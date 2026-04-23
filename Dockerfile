FROM golang:1.22-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/fleet_app ./main.go

FROM alpine:3.20
WORKDIR /app

RUN adduser -D appuser

COPY --from=builder /out/fleet_app /app/fleet_app
COPY migrations /app/migrations
COPY seeds /app/seeds
COPY templates /app/templates
COPY static /app/static

RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

CMD ["/app/fleet_app", "serve"]
