# provenix-examples

Examples, guides, and validation workflows for [Provenix](https://github.com/open-verix/provenix) — a Policy-Driven Software Supply Chain Orchestrator.

## Quick Start

```bash
# Install Provenix (Linux/macOS)
curl -sSL https://raw.githubusercontent.com/open-verix/provenix/main/scripts/install.sh | sh

# Initialize (download vulnerability database)
provenix init

# Attest a container image (keyless signing in CI)
provenix attest nginx:latest

# Attest with a local key (development)
provenix attest nginx:latest --key path/to/key.pem --skip-transparency
```

## Repository Structure

```
provenix-examples/
├── docs/
│   ├── quickstart.md          # 5-minute getting started guide
│   ├── github-actions.md      # GitHub Actions integration guide
│   ├── gitlab-ci.md           # GitLab CI integration guide
│   ├── policies.md            # Policy configuration guide
│   └── vex-workflows.md       # VEX workflow guide
├── examples/
│   ├── github-actions/        # Reusable GitHub Actions workflows
│   │   ├── docker-image.yml
│   │   ├── go-application.yml
│   │   ├── multi-arch-build.yml
│   │   └── policy-gate.yml
│   ├── gitlab-ci/             # GitLab CI pipeline examples
│   │   ├── docker-image.gitlab-ci.yml
│   │   ├── go-application.gitlab-ci.yml
│   │   └── policy-enforcement.gitlab-ci.yml
│   ├── policies/              # provenix.yaml policy examples
│   │   ├── provenix.yaml      # Default policy
│   │   └── provenix-cel.yaml  # CEL custom policy
│   └── vex/                   # VEX workflow examples
└── .github/workflows/
    └── validate-examples.yml  # CI: validates examples on each platform
```

## Platform Support

| Platform | Architecture | Status |
|----------|-------------|--------|
| Linux | amd64 | ✅ Supported |
| Linux | arm64 | ✅ Supported |
| macOS | arm64 (Apple Silicon) | ✅ Supported |
| macOS | amd64 (Intel) | ⚠️ Use Rosetta 2 |
| Windows | amd64 | ✅ Supported |

## Documentation

- [Quick Start](docs/quickstart.md)
- [GitHub Actions Guide](docs/github-actions.md)
- [GitLab CI Guide](docs/gitlab-ci.md)
- [Policy Configuration](docs/policies.md)
- [VEX Workflows](docs/vex-workflows.md)

## Links

- [Provenix Source](https://github.com/open-verix/provenix)
- [Releases](https://github.com/open-verix/provenix/releases)
- [Design Docs](https://github.com/open-verix/provenix/tree/main/docs)
