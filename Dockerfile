# syntax=docker/dockerfile:1.7

# 1. Budowanie aplikacji w Go
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build
WORKDIR /src

# Certyfikaty do HTTPS, a UPX do zmnieszenia rozmiaru
RUN apk add --no-cache ca-certificates upx

# Kopiowanie go.mod, by Docker mógł lepiej używać cache
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Kopiowanie kodu aplikacji
COPY main.go ./

# Zmienne potrzebne przy budowaniu obrazów na różne platformy
ARG TARGETOS
ARG TARGETARCH

# Budowanie statycznej i odchudzonej wersji binarnej aplikacji
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build \
    -trimpath \
    -buildvcs=false \
    -tags="netgo osusergo" \
    -ldflags="-s -w -buildid=" \
    -o /out/weatherapp .

# Kompresowanie binarki, by obraz końcowy był mniejszy
RUN upx --best --lzma /out/weatherapp

# 2. Finalny obraz bez systemu bazowego
FROM scratch

# Dane zgodne z OCI
LABEL org.opencontainers.image.authors="PAWEŁ LADOWSKI"
LABEL org.opencontainers.image.title="Weather App - Zadanie 1"
LABEL org.opencontainers.image.description="Minimalny obraz aplikacji pogodowej"

# Kopiowanie certyfikatów i gotowej aplikacji
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/weatherapp /weatherapp

# Domyślny port aplikacji
ENV PORT=8080
EXPOSE 8080

# Użytkownik non-root dla bezpieczeństwa
USER 65532:65532

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/weatherapp", "healthcheck"]

# Uruchomienie aplikacji
ENTRYPOINT ["/weatherapp"]
