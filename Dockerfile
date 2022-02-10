FROM golang:alpine as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o change-monitor -ldflags="-w -s" cmd/sa.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/change-monitor /usr/bin/

ENTRYPOINT ["change-monitor"]
