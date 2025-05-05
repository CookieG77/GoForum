# syntax=docker/dockerfile:1
################################################################################

# Creation de l'image de build
ARG GO_VERSION=1.24.1
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /app

# Copie des fichiers sources
COPY projet/ .
# Recuperation des dependances
RUN go mod download -x

# Choix de l'OS et de l'architecture
ARG TARGETARCH=amd64

# # Compilation du programme
# RUN CGO_ENABLED=1 GOOS=linux GOARCH=${TARGETARCH} go build -o -v /bin/projet
RUN CGO_ENABLED=1 GOOS=linux GOARCH=${TARGETARCH} go build \
    -ldflags="-s -w" -trimpath -v -o /bin/projet .


# # Creation de l'image Linux
FROM ubuntu:25.04

RUN --mount=type=cache,target=/var/apt-cache \
    apt-get update && apt-get install -y \
        ca-certificates \
        tzdata \
        sqlite3 \
        && \
        update-ca-certificates \
        && apt-get clean \
        && rm -rf /var/lib/apt/lists/*

# Copie du binaire dans l'image finale
COPY --from=build /bin/projet /projet

ARG UID=10001
RUN useradd \
    --no-create-home \
    --home-dir "/nonexistent" \
    --shell "/usr/sbin/nologin" \
    --uid "${UID}" \
    --user-group \
    appuser
RUN mkdir -p /data && chown appuser:appuser /data

# USER appuser
USER appuser

# Ouverture des ports
EXPOSE 8080
EXPOSE 443

# COPY ./docker-entrypoint.sh .
# RUN chmod +x docker-entrypoint.sh
# ENTRYPOINT [ "docker-entrypoint.sh" ]

# RUN ls -R projet || echo "File not found but continuing anyway"

# RUN chmod +x /projet

ENTRYPOINT [ "/projet" ]