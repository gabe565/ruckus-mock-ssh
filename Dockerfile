ARG GO_VERSION=1.19

FROM golang:$GO_VERSION-alpine as go-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache \
    go build -ldflags="-w -s"


FROM alpine
LABEL org.opencontainers.image.source="https://github.com/gabe565/ruckus-mock-ssh"
WORKDIR /app

COPY --from=go-builder /app/ruckus-mock-ssh .

CMD ["./ruckus-mock-ssh", "--address=:2222"]
