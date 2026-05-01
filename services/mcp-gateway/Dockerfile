# Build stage
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/mcp-gateway ./cmd/mcp-gateway

# Runtime stage — distroless, no shell, no package manager
FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/mcp-gateway /usr/local/bin/mcp-gateway
USER nonroot:nonroot
EXPOSE 8443
ENTRYPOINT ["/usr/local/bin/mcp-gateway"]
