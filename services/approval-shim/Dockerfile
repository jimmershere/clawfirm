FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/approval-shim ./cmd/approval-shim

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/approval-shim /usr/local/bin/approval-shim
USER nonroot:nonroot
EXPOSE 50051
ENTRYPOINT ["/usr/local/bin/approval-shim"]
