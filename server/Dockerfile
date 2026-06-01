FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /jafa-server ./cmd/main.go
# dbmate runs DB migrations on container start (see entrypoint.sh)
RUN CGO_ENABLED=0 GOOS=linux go install github.com/amacneil/dbmate/v2@v2.21.0

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /jafa-server .
COPY --from=builder /go/bin/dbmate /usr/local/bin/dbmate
COPY internal/infrastructure/psql/migrations ./migrations
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["./entrypoint.sh"]
