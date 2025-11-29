############################
# Builder
############################
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Ensure static binary for scratch
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN go mod tidy
RUN go build -trimpath -ldflags "-s -w" -o /server ./cmd/api

############################
# Runtime
############################
FROM scratch AS runtime

ENV PORT=8080
EXPOSE 8080

# Run as non-root user
USER 65532:65532

COPY --from=builder /server /server

ENTRYPOINT ["/server"]
