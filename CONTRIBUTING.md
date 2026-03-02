# Contributing

This project wraps [adobe/ims-go](https://github.com/adobe/ims-go).

## Development

```bash
# Build
go build -o imscli .

# Run all tests
go test ./...

# Run a single test by name
go test ./ims/ -run TestValidateURL

# Vet (static analysis)
go vet ./...
```

## CI Pipelines

All pipelines are defined in `.github/workflows/`.

### CI (`ci.yml`)

**Triggers:** Every push to any branch and pull requests targeting `main`.

Runs three parallel jobs:

- **Test** — Runs `go test -race` with coverage and prints a coverage summary to the log. Verifies dependency checksums with `go mod verify`.
- **Lint** — Runs `go vet` and [golangci-lint](https://golangci-lint.run/) for extended static analysis.
- **Build** — Verifies the project compiles and validates `.goreleaser.yml` with `goreleaser check`. Cross-platform builds are handled by GoReleaser at release time, so a single-platform check suffices here.

### PR Title (`pr-title.yml`)

**Triggers:** When a pull request is opened, edited, synchronized or reopened.

Validates that the PR title follows the [Conventional Commits](https://www.conventionalcommits.org) format (e.g. `feat: add token refresh`, `fix(auth): handle expired tokens`). This is enforced because the repository is configured for **squash merging only** — the PR title becomes the commit message on `main`, and GoReleaser uses these prefixes to group the release changelog into Features, Bug Fixes, etc.

### Release (`main.yml`)

**Triggers:** Pushing a version tag (e.g. `v1.2.3`) or manual dispatch from the Actions tab.

Runs [GoReleaser](https://goreleaser.com) to build cross-platform binaries (linux/darwin/windows × amd64/arm64), package them as archives and system packages (deb, rpm, apk), generate a changelog grouped by commit type, and publish a GitHub Release with all artifacts.

### CodeQL (`codeql-analysis.yml`)

**Triggers:** Every push, pull requests targeting `main`, and weekly on a cron schedule.

Runs GitHub's CodeQL security analysis to detect vulnerabilities in the Go source code.

### govulncheck (`govulncheck.yml`)

**Triggers:** Every push to `main`, pull requests, and weekly on Monday at 9:00 UTC.

Runs Go's official vulnerability scanner ([govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)) against all packages. The job fails if any known vulnerabilities in the Go vulnerability database affect the code. Unlike general-purpose scanners, govulncheck traces call graphs — it only reports vulnerabilities in functions your code actually calls. The weekly schedule catches new vulnerabilities even when the code hasn't changed. If the scheduled scan fails, a GitHub issue labeled `security` is created automatically.

## Repository Settings

- **Squash merge only** — Merge commits and rebase merging are disabled. The PR title is used as the squash commit message, ensuring conventional commit messages land on `main`.
- **Auto-delete branches** — Head branches are automatically deleted after a PR is merged.
- **Renovate auto-merge** — [Renovate](https://docs.renovatebot.com/) monitors `go.mod` for dependency updates and opens PRs automatically. Patch updates (e.g., `v1.8.0` → `v1.8.1`) are auto-merged after CI passes. Minor and major updates require manual review. Configuration lives in `renovate.json`.
- **Go version pinning** — All CI workflows use `go-version-file: go.mod` so the Go compiler version is controlled by the `toolchain` directive in `go.mod`. To upgrade Go, update `go.mod` (Renovate opens PRs for this automatically). To downgrade after a bad release, revert the `toolchain` line in `go.mod` — all workflows pick up the change immediately.

## Release Process

In order to standardize the release process, [goreleaser](https://goreleaser.com) has been adopted.

To build and release a new version:
```
git tag vX.X.X && git push --tags
goreleaser release --clean
```

The binary version is set automatically to the git tag.

Please tag new versions using [semantic versioning](https://semver.org/spec/v2.0.0.html).

## Development Notes

### PersistentPreRunE and subcommands

The root command defines a `PersistentPreRunE` that loads configuration from flags, environment variables, and config files (see `cmd/root.go`). In cobra, if a subcommand defines its own `PersistentPreRunE`, it **overrides** the parent's — the root's `PersistentPreRunE` will not run for that subcommand or its children. If you need to add a `PersistentPreRunE` to a subcommand, you must explicitly call the parent's first.

## Additional Reading

The `docs/` directory contains write-ups on non-obvious problems encountered during development:

- [OAuth Local Server: Unhandled Serve() Error](docs/oauth-serve-error.md) — Why the `Serve` goroutine error is captured via a buffered channel, and why the channel must be buffered.
- [OAuth Local Server Shutdown Deadlock](docs/oauth-shutdown-deadlock.md) — How an unbuffered channel creates a deadlock between the HTTP handler and `Shutdown()`, and the two-part fix.
