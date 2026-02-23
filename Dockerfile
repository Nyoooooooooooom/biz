FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY go.mod .
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/biz ./cmd/biz

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends chromium fonts-dejavu ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /out/biz /usr/local/bin/biz
COPY templates ./templates
COPY config.example.yaml ./config.example.yaml
ENTRYPOINT ["/usr/local/bin/biz"]
