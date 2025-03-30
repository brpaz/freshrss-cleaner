# =========================================
# Build stage
# =========================================
FROM --platform=$BUILDPLATFORM golang:1.23-alpine3.21 as build

ARG TARGETOS
ARG TARGETARCH
ARG BUILD_DATE
ARG GIT_COMMIT
ARG GIT_VERSION

WORKDIR /app

# Copy only dependency files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies as a separate layer
RUN go mod download && go mod verify

# Build with optimizations and security considerations
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build \
    -trimpath \
    -ldflags="-s -w \
    -X github.com/brpaz/freshrss-cleaner/cmd/version.BuildDate=${BUILD_DATE} \
    -X github.com/brpaz/freshrss-cleaner/cmd/version.Version=${GIT_VERSION} \
    -X github.com/brpaz/freshrss-cleaner/cmd/version.GitCommit=${GIT_COMMIT} \
    -extldflags '-static'" \
    -o /out/freshrss-cleaner ./main.go

# ====================================
# Production stage
# ====================================
FROM alpine:3.21 as prod

ENV PUID=1000
ENV PGID=1000

# Add CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata su-exec shadow && \
    update-ca-certificates

# Create a non-root user and group
RUN addgroup -g ${PGID} app && \
    adduser -D -u ${PUID} -G app app && \
    mkdir -p /app && \
    chown app:app /app

# Copy entrypoint script
COPY docker/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Copy binary from build stage
COPY --from=build --chown=app:app /out/freshrss-cleaner /app/

WORKDIR /app

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]

CMD ["/app/freshrss-cleaner"]

