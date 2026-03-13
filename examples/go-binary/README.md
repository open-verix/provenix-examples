# Go Binary Example

This example demonstrates generating atomic evidence for a simple Go CLI application using Provenix.

## Overview

**Artifact Type:** Go Binary (statically-linked executable)

**What This Example Demonstrates:**

- ✅ Binary SBOM generation with Syft
- ✅ Vulnerability scanning with Grype
- ✅ Keyless signing using GitHub OIDC
- ✅ Publishing attestations to Rekor transparency log
- ✅ Policy-based compliance checking
- ✅ Multi-architecture builds

## Quick Start

### Prerequisites

- Go 1.22 or later
- [Provenix](https://github.com/open-verix/provenix/releases) (latest release)

### Install Provenix

```bash
# macOS/Linux
curl -sSL https://raw.githubusercontent.com/open-verix/provenix/main/scripts/install.sh | sh

# Verify
provenix version
```

### Local Development

```bash
# Change directory
cd ~/provenix-examples/examples/go-binary

# Build the binary
make build

# Run the application
make run

# Check version information
make version
```

### Generate Attestation (Local / Air-gapped)

```bash
# Build the binary
make build

# Step 1: Initialize Grype vulnerability database (one-time, ~200MB download)
provenix init

# Step 2 (optional): Generate a key pair for local signing
# Creates cosign.key (private) and cosign.pub (public)
# Note: cosign.key is a private key — do not commit it to version control
provenix init --generate-key --key-output cosign

# Step 3: Generate attestation (skip Rekor for local testing)
provenix attest build/app \
  --key cosign.key \
  --output attestation.json \
  --skip-transparency

# View the attestation
cat attestation.json | jq .
```

### Generate Attestation (GitHub Actions — Keyless)

The recommended approach uses GitHub Actions with OIDC for keyless signing.  
See [.github/workflows/go-binary.yml](../../.github/workflows/go-binary.yml) for the complete workflow.

## Project Structure

```
go-binary/
├── main.go          # Simple CLI application
├── go.mod           # Go module definition
├── Makefile         # Build automation
├── provenix.yaml    # Provenix configuration (tool config + policy: section)
├── policy.yaml      # Deprecated — policy is now in provenix.yaml
└── README.md        # This file
```

## Complete Workflow

### Step 1: Build

```bash
make build
# Output: build/app
```

The binary is built with version information injected via ldflags:

- Version: From git tags
- Commit: From `git rev-parse`
- BuildDate: Current UTC timestamp

### Step 2: Generate Attestation

```bash
# Reads provenix.yaml automatically; use --key for local (non-CI) signing
provenix attest build/app \
  --output attestation.json \
  --skip-transparency
```

### Step 3: Generate VEX (optional)

```bash
provenix vex generate --output vex.json
```

### Step 4: Policy Check

```bash
provenix policy check
```

## Multi-Architecture Builds

> **Note:** `make build-multi` uses inline environment variable syntax (`GOOS=linux go build`)  
> which requires a Unix-compatible shell. Use **WSL** or **Git Bash** on Windows.

```bash
make build-multi
```

Produces binaries for:

- `build/app-linux-amd64`
- `build/app-linux-arm64`
- `build/app-darwin-arm64`
- `build/app-windows-amd64.exe`

## CI/CD Integration

See the GitHub Actions workflow at [go-binary.yml](../../.github/workflows/go-binary.yml):

- Triggers on changes to `examples/go-binary/`
- Builds the binary
- Generates keyless attestation via OIDC
- Uploads attestation as artifact

## Related Resources

- [Provenix CLI Reference](https://github.com/open-verix/provenix/blob/main/docs/drafts/cli_specification.md)
- [Policy Configuration](https://github.com/open-verix/provenix/blob/main/docs/drafts/configuration.md)
- [Back to Examples](../../README.md)
