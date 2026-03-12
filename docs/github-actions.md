# GitHub Actions Integration Guide

Integrate Provenix into your GitHub Actions workflows to automatically generate **atomic evidence** (SBOM + Vulnerability Scan + Signature) on every build.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Step](#installation-step)
- [Workflow Permissions](#workflow-permissions)
- [Examples](#examples)
  - [Go Application](#go-application)
  - [Security Policy Gate](#security-policy-gate)
  - [Multi-Architecture Build](#multi-architecture-build)
  - [Docker Image](#docker-image)
- [Exit Codes](#exit-codes)
- [Tips & Best Practices](#tips--best-practices)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

- GitHub repository with Actions enabled
- Provenix `v0.1.0-alpha.2` or later
- Go `1.22`+ (for Go projects)

---

## Installation Step

Add this step to your workflow once — `install.sh` auto-detects the platform (Linux/macOS, amd64/arm64):

```yaml
- name: Install Provenix
  run: |
    curl -sSL https://raw.githubusercontent.com/open-verix/provenix/main/scripts/install.sh | sh
    provenix version
```

After this step, `provenix` is available in `$PATH` for all subsequent steps.

> **No manual platform detection needed.** `install.sh` selects the correct binary for `ubuntu-latest`, `macos-latest`, etc.

---

## Workflow Permissions

Keyless signing uses GitHub OIDC. Always declare these permissions at the job level:

```yaml
jobs:
  build-and-attest:
    permissions:
      id-token: write # Required for keyless signing (OIDC)
      contents: write # Required if uploading to releases
```

For **organization repositories**, also enable in:
`Settings → Actions → General → Workflow permissions → Allow GitHub Actions to create and approve pull requests`

---

## Examples

### Go Application

The simplest workflow: build → attest → upload.

**Full template:** [`examples/github-actions/go-application.yml`](../examples/github-actions/go-application.yml)  
**Working CI reference:** [`.github/workflows/go-binary.yml`](../.github/workflows/go-binary.yml)

```yaml
name: Attest Go Application

on:
  push:
    branches: [main]
    tags: ["v*"]

jobs:
  build-and-attest:
    runs-on: ubuntu-latest
    permissions:
      id-token: write
      contents: write

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: Build
        run: |
          go build -trimpath \
            -ldflags="-s -w -X main.version=${{ github.ref_name }}" \
            -o myapp ./cmd/myapp

      - name: Install Provenix
        run: curl -sSL https://raw.githubusercontent.com/open-verix/provenix/main/scripts/install.sh | sh

      - name: Initialize
        run: provenix init

      - name: Attest
        run: provenix attest myapp --output attestation.json

      - uses: actions/upload-artifact@v4
        with:
          name: attestation
          path: attestation.json
```

---

### Security Policy Gate

Fail the build when vulnerabilities exceed policy thresholds. Use this as a PR check.

**Full template:** [`examples/github-actions/policy-gate.yml`](../examples/github-actions/policy-gate.yml)  
**Policy file reference:** [`examples/go-binary/policy.yaml`](../examples/go-binary/policy.yaml)

```yaml
- name: Attest (no Rekor for PRs)
  run: |
    provenix attest myapp \
      --output attestation.json \
      --skip-transparency

# exit 0 = pass, exit 1 = policy violation
- name: Policy check
  run: provenix policy check --config provenix.yaml
```

**Policy file (`provenix.yaml`):**

```yaml
policy:
  max_critical: 0
  max_high: 5
  custom:
    - name: no-critical-vulns
      expression: "vulnerabilities.critical == 0"
```

> **Exit code semantics:**
>
> - `0` — policy passed
> - `1` — policy violation (build should fail)
> - `2` — attestation saved locally, Rekor unavailable (partial success)

---

### Multi-Architecture Build

Build and attest for multiple platforms in a matrix.

**Full template:** [`examples/github-actions/multi-arch-build.yml`](../examples/github-actions/multi-arch-build.yml)

```yaml
jobs:
  build-matrix:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest]
        arch: [amd64, arm64]
        exclude:
          - os: ubuntu-latest
            arch: arm64 # Cross-compile instead

    runs-on: ${{ matrix.os }}
    permissions:
      id-token: write

    steps:
      - name: Install Provenix
        # install.sh auto-detects OS and arch
        run: curl -sSL https://raw.githubusercontent.com/open-verix/provenix/main/scripts/install.sh | sh

      - name: Attest
        run: |
          provenix attest myapp-${{ matrix.os }}-${{ matrix.arch }} \
            --output attestation-${{ matrix.os }}-${{ matrix.arch }}.json
```

---

### Docker Image

Attest a Docker image by its digest (prevents TOCTOU attacks).

**Full template:** [`examples/github-actions/docker-image.yml`](../examples/github-actions/docker-image.yml)

```yaml
- name: Build and push Docker image
  id: build
  uses: docker/build-push-action@v5
  with:
    push: true
    tags: ${{ steps.meta.outputs.tags }}

- name: Attest by digest (not tag)
  run: |
    IMAGE="${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}@${{ steps.build.outputs.digest }}"
    provenix attest "${IMAGE}" --output attestation.json
```

> **Always attest by digest**, not by tag. Tags are mutable — digests are not.

---

## Exit Codes

| Code | Meaning                                            | CI Behavior  |
| ---- | -------------------------------------------------- | ------------ |
| `0`  | Complete success (signed + published to Rekor)     | Pass         |
| `1`  | Fatal error (crypto failure, artifact not found)   | Fail         |
| `2`  | Partial success (saved locally, Rekor unavailable) | Configurable |

To allow exit code `2` (e.g., air-gapped or Rekor outage):

```yaml
- name: Attest
  run: provenix attest myapp --output attestation.json || [ $? -eq 2 ]
```

---

## Tips & Best Practices

### Use `--skip-transparency` for Pull Requests

Rekor publishes attestations permanently. For PRs, skip publishing to avoid noise:

```yaml
- name: Attest (PR)
  if: github.event_name == 'pull_request'
  run: provenix attest myapp --output attestation.json --skip-transparency

- name: Attest (main/tags)
  if: github.event_name != 'pull_request'
  run: provenix attest myapp --output attestation.json
```

### Cache the Grype Database

`provenix init` downloads a ~200MB vulnerability database. Cache it between runs:

```yaml
- name: Cache Grype DB
  uses: actions/cache@v4
  with:
    path: ~/.cache/grype
    key: grype-db-${{ runner.os }}-${{ hashFiles('**/grype.db') }}
    restore-keys: grype-db-${{ runner.os }}-

- name: Initialize
  run: provenix init
```

### Use a `provenix.yaml` configuration file

Commit a `provenix.yaml` to your repository to codify your SBOM format and scan policy:

```yaml
sbom:
  format: cyclonedx-json

scan:
  min-severity: medium

signing:
  keyless: true
  rekor: true
```

Then reference it with `--config provenix.yaml` in your workflow.

---

## Troubleshooting

### `Error: failed to get OIDC token`

**Cause:** `id-token: write` permission is missing.  
**Fix:** Add `id-token: write` to your job's `permissions` block.

### `Error: failed to publish to Rekor`

**Cause:** Network connectivity issue or Rekor outage.  
**Fix:** This exits with code `2`. The attestation is saved locally. Use `--skip-transparency` if Rekor is unavailable.

### `provenix: command not found`

**Cause:** `install.sh` path not added to `$PATH`, or step order issue.  
**Fix:** Ensure the Install step runs before any `provenix` command. Check the shell output — `install.sh` prints the install path.

### Policy check fails with unknown vulnerabilities

**Cause:** Grype database is outdated.  
**Fix:** Run `provenix init` before `provenix attest` to refresh the database.
