# Quick Start Guide

Get Provenix running in 5 minutes.

## 1. Install

### Linux / macOS (recommended)

```bash
curl -sSL https://raw.githubusercontent.com/open-verix/provenix/main/scripts/install.sh | sh
```

### Manual download

Download the latest binary from [GitHub Releases](https://github.com/open-verix/provenix/releases):

| Platform | File |
|----------|------|
| Linux amd64 | `provenix_<version>_linux_amd64.tar.gz` |
| Linux arm64 | `provenix_<version>_linux_arm64.tar.gz` |
| macOS arm64 | `provenix_<version>_darwin_arm64.tar.gz` |
| Windows amd64 | `provenix_<version>_windows_amd64.zip` |

```bash
# Example: Linux amd64
VERSION=v0.1.0-alpha.2
curl -sSL "https://github.com/open-verix/provenix/releases/download/${VERSION}/provenix_${VERSION}_linux_amd64.tar.gz" | tar xz
chmod +x provenix
sudo mv provenix /usr/local/bin/
```

### Build from source

```bash
git clone https://github.com/open-verix/provenix.git
cd provenix
go build -o provenix ./cmd/provenix
```

## 2. Initialize

Download the Grype vulnerability database (~200MB, one-time):

```bash
provenix init
```

## 3. First attestation

### With keyless signing (requires OIDC — recommended in CI)

```bash
provenix attest nginx:latest
```

### With local key (development / air-gapped)

```bash
# Generate a key pair
openssl ecparam -name prime256v1 -genkey -noout -out key.pem
openssl ec -in key.pem -pubout -out pub.pem

# Attest
provenix attest nginx:latest --key key.pem --skip-transparency
```

## 4. Verify

```bash
provenix verify attestation.json --key pub.pem
```

## 5. Generate a report

```bash
provenix report attestation.json --format markdown
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Complete success (signed + published to Rekor) |
| 1 | Fatal error |
| 2 | Partial success (saved locally, Rekor unavailable) |

## Next Steps

- [GitHub Actions integration](github-actions.md)
- [Policy configuration](policies.md)
- [VEX workflows](vex-workflows.md)
