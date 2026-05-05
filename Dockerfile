# ---- Build ----
FROM golang:1.22-alpine AS build
WORKDIR /app

# Cache the module graph first.
COPY go.mod go.sum ./
RUN go mod download

# Compile the static, embedded-asset binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /server ./cmd/server

# ---- Runtime ----
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=build /server ./server
COPY configs/ ./configs/

EXPOSE 8080

ENTRYPOINT ["./server"]
CMD ["serve"]
