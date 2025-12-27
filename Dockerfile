FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o ./build/ ./cmd/...

FROM alpine
WORKDIR /app
COPY --from=builder /app/build/api ./api
CMD ["/app/api"]
