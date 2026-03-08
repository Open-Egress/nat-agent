# Open Egress Agent 🛰️

The **Open Egress Agent** is the lightweight data-plane component of the Open Egress ecosystem. It is a stateless Go binary designed to run on cloud instances (EC2, VM, Compute Instance) to transform them into high-performance, self-healing NAT Gateways.

## 🛠 Key Functions

- **Kernel Orchestration**: Automatically enables IPv4 forwarding and configures nftables/iptables masquerading.
- **Health Monitoring**: Provides a constant heartbeat to the Open Egress Control Plane.
- **Metrics Collection**: Scrapes interface statistics (`/proc/net/dev`) to track egress volume without packet inspection overhead.
- **Zero-Dependency**: Compiled as a static binary with no CGO dependencies—runs on any Linux distro.

## 🏗 Repository Structure

```text
.
├── cmd/agent/                # Entry point & Signal handling
├── internal/
│   ├── heartbeat/            # Logic for phoning home to Control Plane
│   ├── kernel/               # Kernel parameter tuning (sysctl)
│   ├── routing/              # nftables/iptables rule management
│   └── stats/                # Low-overhead traffic counters
├── Makefile                  # Multi-arch build system (ARM64/AMD64)
└── open-egress-agent.service # Systemd unit for auto-restart
```

## 🚀 Getting Started

### 1. Build from Source

Ensure you have **Go 1.24+** installed. Use the provided Makefile to build for your architecture (use ARM64 for the best cost-to-performance ratio).

```bash
# Build for AWS Graviton or Oracle Ampere
make build-arm64

# Build for Intel/AMD
make build-amd64
```

### 2. Configuration

The agent is configured via environment variables or a `.env` file:

| Variable      | Description                               | Example                           |
|---------------|-------------------------------------------|-----------------------------------|
| `CONTROL_URL` | URL of your Open Egress Control Plane     | `https://control.egress.local`    |
| `API_KEY`   | Auth key for the Control Plane            | `your-secure-token`               |
| `INTERFACE`   | The public-facing network interface       | `eth0`                            |
| `HB_INTERVAL` | Heartbeat frequency                       | `10s`                             |

### 3. Deployment (Systemd)

To ensure the agent stays running and starts on boot:

```bash
sudo cp dist/open-egress-agent-linux-arm64 /usr/local/bin/open-egress-agent
sudo cp open-egress-agent.service /etc/systemd/system/
sudo systemctl enable --now open-egress-agent
```

## ⚡ Performance Tuning

The agent automatically optimizes the Linux networking stack for high-concurrency NAT by tuning the following (optional/configurable):

- `net.ipv4.ip_forward = 1`
- `net.netfilter.nf_conntrack_max` (increased for high-traffic nodes)
- `net.ipv4.tcp_fin_timeout = 15`

## 📊 Logging & Observability

The agent uses structured logging with Go's `log/slog` for better observability and easier parsing in cloud environments. Logs include context such as:
- **Operations**: `Adding nftables table`, `Kernel parameter updated successfully`.
- **Context**: `name`, `family`, `path`, `value`, `error`.

Example log entry:
`2026/03/08 17:03:21 INFO Adding nftables table name=nat family=1`

## 📄 License

Licensed under the Apache License 2.0.
