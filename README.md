# vaultwatch

A daemon that monitors HashiCorp Vault secret lease expirations and triggers configurable webhook alerts.

---

## Installation

```bash
go install github.com/yourusername/vaultwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultwatch.git && cd vaultwatch && go build -o vaultwatch .
```

---

## Usage

Create a configuration file (`config.yaml`):

```yaml
vault:
  address: "https://vault.example.com"
  token: "s.yourVaultToken"

watch_interval: "5m"
alert_threshold: "24h"

webhooks:
  - name: "slack-alerts"
    url: "https://hooks.slack.com/services/your/webhook/url"
    on_events: ["lease_expiring", "lease_expired"]
```

Run the daemon:

```bash
vaultwatch --config config.yaml
```

vaultwatch will poll Vault at the configured interval and fire webhook payloads when a secret lease is within the alert threshold or has already expired.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to configuration file |
| `--log-level` | `info` | Log verbosity (`debug`, `info`, `warn`, `error`) |
| `--dry-run` | `false` | Detect expirations without firing webhooks |

---

## Requirements

- Go 1.21+
- HashiCorp Vault with a token that has `read` access to monitored secret paths

---

## License

MIT © 2024 yourusername