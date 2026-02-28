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

- **Test** — Runs `go test -race` with coverage and uploads the report to Codecov.
- **Lint** — Runs `go vet` and [golangci-lint](https://golangci-lint.run/) for extended static analysis.
- **Build** — Verifies the project compiles. Cross-platform builds are handled by GoReleaser at release time, so a single-platform check suffices here.

### PR Title (`pr-title.yml`)

**Triggers:** When a pull request is opened, edited, synchronized or reopened.

Validates that the PR title follows the [Conventional Commits](https://www.conventionalcommits.org) format (e.g. `feat: add token refresh`, `fix(auth): handle expired tokens`). This is enforced because the repository is configured for **squash merging only** — the PR title becomes the commit message on `main`, and GoReleaser uses these prefixes to group the release changelog into Features, Bug Fixes, etc.

### Release (`main.yml`)

**Triggers:** Pushing a version tag (e.g. `v1.2.3`) or manual dispatch from the Actions tab.

Runs [GoReleaser](https://goreleaser.com) to build cross-platform binaries (linux/darwin/windows × amd64/arm64), package them as archives and system packages (deb, rpm, apk), generate a changelog grouped by commit type, and publish a GitHub Release with all artifacts.

### CodeQL (`codeql-analysis.yml`)

**Triggers:** Every push, pull requests targeting `main`, and weekly on a cron schedule.

Runs GitHub's CodeQL security analysis to detect vulnerabilities in the Go source code.

## Repository Settings

- **Squash merge only** — Merge commits and rebase merging are disabled. The PR title is used as the squash commit message, ensuring conventional commit messages land on `main`.
- **Auto-delete branches** — Head branches are automatically deleted after a PR is merged.

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
