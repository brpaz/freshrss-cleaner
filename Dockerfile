# =========================================
# Build stage
# =========================================
FROM --platform=$BUILDPLATFORM golang:1.24-alpine3.21 AS build

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
FROM alpine:3.21 AS prod

RUN apk add --no-cache ca-certificates curl

# Copy binary from build stage
COPY --from=build /out/freshrss-cleaner /usr/local/bin/freshrss-cleaner

ENTRYPOINT ["/usr/local/bin/freshrss-cleaner"]


