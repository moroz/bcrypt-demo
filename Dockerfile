ARG GO_VERSION=1.21.4
ARG DEBIAN_VERSION=bookworm
ARG NODE_VERSION=20.9.0

FROM golang:${GO_VERSION}-${DEBIAN_VERSION} AS builder

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN templ generate

RUN go build -o server

FROM node:$NODE_VERSION AS assets

WORKDIR /app/assets

RUN npm i -g pnpm

COPY assets/package.json assets/pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY assets/ ./

RUN pnpm build

FROM debian:${DEBIAN_VERSION}-slim AS runner

WORKDIR /app

RUN apt-get update -y && apt-get install -y curl \
  && apt-get clean && rm -f /var/lib/apt/lists/*_*

RUN /bin/bash -c 'set -ex && \
    ARCH=`uname -m` && \
    if [ "$ARCH" == "x86_64" ] || [ "$ARCH" == "amd64" ]; then \
      curl --output migrate.tar.gz -LJO https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz; \
    else \
      curl --output migrate.tar.gz -LJO https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-arm64.tar.gz; \
    fi' && tar xzf migrate.tar.gz -C /usr/local/bin && rm migrate.tar.gz

COPY --from=builder /app/server /app/server
COPY --from=assets /app/static /app/static
COPY db/ /app/db
COPY scripts/ /app/scripts/

CMD ["./scripts/entrypoint.sh"]
