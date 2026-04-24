# vaultdiff

> CLI tool to diff and audit changes between HashiCorp Vault secret versions

---

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultdiff.git
cd vaultdiff
go build -o vaultdiff .
```

---

## Usage

Ensure your Vault environment variables are set (`VAULT_ADDR`, `VAULT_TOKEN`), then run:

```bash
# Diff two versions of a secret
vaultdiff secret/data/myapp/config --v1 3 --v2 4

# Audit all version changes for a secret path
vaultdiff secret/data/myapp/config --audit

# Output diff in JSON format
vaultdiff secret/data/myapp/config --v1 1 --v2 2 --format json
```

### Example Output

```diff
secret/data/myapp/config (v3 → v4)

  DB_HOST:     "db.internal"
- DB_PORT:     "5432"
+ DB_PORT:     "5433"
- API_KEY:     "old-key-abc123"
+ API_KEY:     "new-key-xyz789"
```

---

## Requirements

- Go 1.21+
- HashiCorp Vault with KV v2 secrets engine enabled

---

## License

[MIT](LICENSE)