# Simple client to test connectivity behind the NAT
FROM alpine:latest

# Install tools for networking tests
RUN apk add --no-cache curl iproute2 traceroute mtr

# Keep the container running
CMD ["tail", "-f", "/dev/null"]
