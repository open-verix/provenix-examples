# Provenix Examples

**Comprehensive examples demonstrating Provenix attestation workflows across different environments**

[![Validate Examples](https://github.com/open-verix/provenix-examples/actions/workflows/validate-examples.yml/badge.svg)](https://github.com/open-verix/provenix-examples/actions/workflows/validate-examples.yml)
[![Go Binary](https://github.com/open-verix/provenix-examples/actions/workflows/go-binary.yml/badge.svg)](https://github.com/open-verix/provenix-examples/actions/workflows/go-binary.yml)
[![Validate Templates](https://github.com/open-verix/provenix-examples/actions/workflows/validate-templates.yml/badge.svg)](https://github.com/open-verix/provenix-examples/actions/workflows/validate-templates.yml)

---

## 🎯 Purpose

This repository provides real-world examples of using [Provenix](https://github.com/open-verix/provenix) to generate **atomic evidence** (SBOM + Vulnerability Scan + Signature) for various types of software artifacts.

**What is Atomic Evidence?**
A cryptographically-signed package containing:

- 📦 **Subject**: The artifact being attested (binary, container, library, etc.)
- 📋 **SBOM**: Software Bill of Materials (CycloneDX or SPDX)
- 🔍 **Vulnerability Report**: Security scan results (Grype)
- ✍️ **Signature**: Cryptographic proof of integrity (Cosign)

All generated **atomically** to prevent TOCTOU (Time-of-Check to Time-of-Use) vulnerabilities.

---

## 📚 Examples

| Example                          | Artifact Type | Language | Highlights                           |
| -------------------------------- | ------------- | -------- | ------------------------------------ |
| [go-binary](examples/go-binary/) | Binary        | Go       | Statically-linked, multi-arch builds |
| docker-image _(coming soon)_     | Container     | Multi    | Multi-stage builds, distroless       |
| nodejs-library _(coming soon)_   | Library       | Node.js  | npm packages, lockfile SBOM          |
| python-package _(coming soon)_   | Library       | Python   | PyPI packages, wheel files           |
| monorepo-app _(coming soon)_     | Mixed         | Multi    | Batch processing, multiple artifacts |

---

## � Documentation

| Guide                                    | Description                                          |
| ---------------------------------------- | ---------------------------------------------------- |
| [Quick Start](docs/quickstart.md)        | Install Provenix and generate your first attestation |
| [GitHub Actions](docs/github-actions.md) | Integrate Provenix into GitHub Actions workflows     |
| GitLab CI _(coming soon)_                | Integrate with GitLab CI/CD pipelines                |
| Policies _(coming soon)_                 | Write and apply security policies                    |
| VEX Workflows _(coming soon)_            | Triage and suppress false-positive vulnerabilities   |

---

## �🚀 Quick Start

### Prerequisites

1. **Install Provenix**

   ```bash
   # macOS/Linux (recommended)
   curl -sSL https://raw.githubusercontent.com/open-verix/provenix/main/scripts/install.sh | sh

   # Verify installation
   provenix version
   ```

   Or download manually from [GitHub Releases](https://github.com/open-verix/provenix/releases):

   | Platform      | File                                     |
   | ------------- | ---------------------------------------- |
   | Linux amd64   | `provenix_<version>_linux_amd64.tar.gz`  |
   | Linux arm64   | `provenix_<version>_linux_arm64.tar.gz`  |
   | macOS arm64   | `provenix_<version>_darwin_arm64.tar.gz` |
   | Windows amd64 | `provenix_<version>_windows_amd64.zip`   |

2. **Clone this repository**

   ```bash
   git clone https://github.com/open-verix/provenix-examples.git
   cd provenix-examples
   ```

### Try Your First Example

```bash
# Go to the simplest example
cd examples/go-binary

# Build the binary
make build

# Generate attestation (local development mode)
provenix attest build/app \
  --output attestation.json \
  --skip-transparency

# View the attestation
cat attestation.json | jq .

# Check policy compliance
provenix policy check attestation.json \
  --policy policy.yaml
```

---

## 🏗️ Repository Structure

```
provenix-examples/
├── examples/
│   ├── go-binary/             # 🔹 Start here
│   │   ├── main.go
│   │   ├── go.mod
│   │   ├── Makefile
│   │   ├── provenix.yaml
│   │   ├── policy.yaml
│   │   └── README.md
│   ├── github-actions/        # Reusable GitHub Actions workflows
│   ├── gitlab-ci/             # GitLab CI pipeline examples
│   ├── policies/              # provenix.yaml policy examples
│   └── vex/                   # VEX workflow examples
├── docs/
│   ├── quickstart.md
│   ├── github-actions.md
│   ├── gitlab-ci.md
│   ├── policies.md
│   └── vex-workflows.md
└── .github/workflows/
    ├── validate-examples.yml  # Multi-platform binary validation
    └── go-binary.yml          # Go binary example CI
```

---

## 🔄 Common Workflows

### 1. Generate Attestation

```bash
# Basic attestation (local key)
provenix attest <artifact> \
  --output attestation.json \
  --key cosign.key

# Keyless attestation (GitHub Actions OIDC)
provenix attest <artifact> \
  --output attestation.json
```

### 2. Generate VEX (Vulnerability Exploitability eXchange)

```bash
provenix vex generate attestation.json \
  --output vex.json \
  --status not_affected \
  --justification "vulnerable_code_not_in_execute_path"
```

### 3. Policy-Based Compliance

```bash
provenix policy check attestation.json \
  --policy examples/policies/provenix.yaml
```

### 4. Batch Processing

```bash
provenix batch \
  --input batch-config.json \
  --parallel 4 \
  --output-dir attestations/
```

---

## 🔐 Keyless Signing with GitHub Actions

All examples include GitHub Actions workflows demonstrating **keyless signing** using OIDC:

```yaml
permissions:
  id-token: write # Required for OIDC token

steps:
  - name: Generate Attestation
    run: |
      provenix attest myapp \
        --output attestation.json
```

**Benefits:**

- ✅ No secret key management
- ✅ Short-lived certificates (automatically rotated)
- ✅ Transparency log (Rekor) ensures tamper-proofing
- ✅ Identity binding to CI/CD workflow

---

## 🖥️ Platform Support

| Platform | Architecture          | Status           |
| -------- | --------------------- | ---------------- |
| Linux    | amd64                 | ✅ Supported     |
| Linux    | arm64                 | ✅ Supported     |
| macOS    | arm64 (Apple Silicon) | ✅ Supported     |
| macOS    | amd64 (Intel)         | ⚠️ Use Rosetta 2 |
| Windows  | amd64                 | ✅ Supported     |

---

## 📝 Policy Examples

### Default Policy (Permissive)

```yaml
version: v1
vulnerabilities:
  max_critical: 0
  max_high: 5
```

### CEL Custom Policy

```yaml
version: v1
custom:
  cel_enabled: true
  cel_expressions:
    - name: no-critical
      expr: "vulnerabilities.filter(v, v.severity == 'Critical').size() == 0"
```

See [examples/policies/](examples/policies/) for complete examples.

---

## 🔗 Additional Resources

- [Provenix Source](https://github.com/open-verix/provenix)
- [Releases](https://github.com/open-verix/provenix/releases)
- [Design Docs](https://github.com/open-verix/provenix/tree/main/docs)
- [CLI Reference](https://github.com/open-verix/provenix/blob/main/docs/drafts/cli_specification.md)
- [in-toto Specification](https://github.com/in-toto/attestation)
- [Sigstore Project](https://www.sigstore.dev/)
- [SLSA Framework](https://slsa.dev/)

---

## 💬 Support

- **Issues**: [GitHub Issues](https://github.com/open-verix/provenix-examples/issues)
- **Discussions**: [GitHub Discussions](https://github.com/open-verix/provenix-examples/discussions)
- **Main Project**: [Provenix Repository](https://github.com/open-verix/provenix)
