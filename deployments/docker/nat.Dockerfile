FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -o open-egress-agent ./cmd/agent


FROM ubuntu:22.04

RUN apt update && \
  apt install -y nftables iproute2 conntrack curl tcpdump && \
  apt clean

COPY --from=builder /app/open-egress-agent /usr/local/bin/open-egress-agent

CMD ["/usr/local/bin/open-egress-agent"]
