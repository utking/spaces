FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN ls -Fla && \
    CGO_ENABLED=0 go test -v ./... -cover && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-extldflags=-static -s -w -extldflags=-static" -o bin/runner cmd/main.go


FROM gcr.io/distroless/static:nonroot

USER nonroot:nonroot
WORKDIR /app

COPY --from=builder /app/bin/runner /app/runner

ENTRYPOINT ["/app/runner", "serve"]
