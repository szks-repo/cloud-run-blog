# syntax=docker/dockerfile:1.6

ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-bookworm AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -buildvcs=false \
    -ldflags "-s -w" \
    -o /cloud-run-blog ./cmd/api

FROM gcr.io/distroless/base-debian12:nonroot AS runtime

ENV PORT=8080
USER nonroot:nonroot
COPY --from=build /cloud-run-blog /cloud-run-blog

EXPOSE 8080
ENTRYPOINT ["/cloud-run-blog"]
