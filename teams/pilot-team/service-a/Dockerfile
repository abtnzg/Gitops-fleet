# syntax=docker/dockerfile:1.6
# Multi-stage build for service-a. The build stage is hermetic (no host Go
# needed), and the final image is distroless:nonroot — no shell, no package
# manager, runs as UID 65532. Matches the platform's hardening bar.

# ---- Build stage ---------------------------------------------------------
FROM golang:1.23-alpine AS build
WORKDIR /src

# Layer-cache dependencies separately from source. service-a has no third-
# party deps (stdlib only), so go.mod is minimal — but the structure is in
# place for when real deps land.
COPY src/go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY src/ ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w \
      -X main.version=$(date -u +%Y%m%d-%H%M%S) \
      -X main.commit=local" \
      -o /out/service-a ./

# ---- Final stage ---------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/service-a /service-a
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/service-a"]
