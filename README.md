# imscli
![CodeQL](https://github.com/adobe/imscli/workflows/CodeQL/badge.svg)
[![govulncheck](https://github.com/adobe/imscli/actions/workflows/govulncheck.yml/badge.svg)](https://github.com/adobe/imscli/actions/workflows/govulncheck.yml)
[![CI](https://github.com/adobe/imscli/actions/workflows/ci.yml/badge.svg)](https://github.com/adobe/imscli/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/adobe/imscli)](https://goreportcard.com/report/github.com/adobe/imscli)
[![Release with GoReleaser](https://github.com/adobe/imscli/actions/workflows/main.yml/badge.svg)](https://github.com/adobe/imscli/actions/workflows/main.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/adobe/imscli)
[![Go Reference](https://pkg.go.dev/badge/github.com/adobe/imscli.svg)](https://pkg.go.dev/github.com/adobe/imscli)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

A CLI tool to troubleshoot and automate Adobe IMS integrations.

## Installation

### Prebuilt binaries

Download the latest release for your platform from the [releases page](https://github.com/adobe/imscli/releases).

### From source

```sh
go install github.com/adobe/imscli@latest
```

## Quick Start

```sh
# Authorize as a user (launches browser for OAuth2 flow)
imscli authorize user --clientID <client-id> --clientSecret <secret> --organization <org> --scopes openid

# Validate a token
imscli validate accessToken --clientID <client-id> --accessToken <token>

# Decode a JWT locally (no API call)
imscli decode --token <jwt>
```

## Commands

| Command | Description |
|---------|-------------|
| `authorize user` | OAuth2 Authorization Code Grant Flow (launches browser) |
| `authorize pkce` | Authorization Code Grant Flow with PKCE (mandatory for public clients, optional for private) |
| `authorize service` | Service authorization (client credentials + service token) |
| `authorize jwt` | JWT Bearer Flow (signed JWT exchanged for access token) |
| `authorize client` | Client Credentials Grant Flow |
| `validate` | Validate a token using the IMS API |
| `invalidate` | Invalidate a token using the IMS API |
| `decode` | Decode a JWT token locally |
| `refresh` | Refresh an access token |
| `exchange` | Cluster access token exchange across IMS Orgs |
| `profile` | Retrieve user profile |
| `organizations` | List user organizations |
| `admin` | Admin operations (profile, organizations) via service token |
| `dcr` | Dynamic Client Registration |

See [DOCUMENTATION.md](DOCUMENTATION.md) for full details on each command.

## Global Flags

These flags apply to all commands:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--url` | `-U` | `https://ims-na1.adobelogin.com` | IMS endpoint URL |
| `--proxyUrl` | `-P` | | HTTP(S) proxy (`http(s)://host:port`) |
| `--proxyIgnoreTLS` | `-T` | `false` | Skip TLS verification (proxy only) |
| `--configFile` | `-f` | | Configuration file path |
| `--timeout` | | `30` | HTTP client timeout in seconds |
| `--verbose` | `-v` | `false` | Verbose output |

## Configuration

Parameters can be provided from three sources (highest to lowest priority):

1. **CLI flags** — `imscli authorize user --scopes openid`
2. **Environment variables** — `IMS_SCOPES=openid imscli authorize user`
3. **Configuration file** — `~/.config/imscli.yaml` or specified with `-f`

See [DOCUMENTATION.md](DOCUMENTATION.md) for configuration file format and examples.

## Contributing

Contributions are welcomed! Read the [Contributing Guide](CONTRIBUTING.md) for more information.

## Licensing

This project is licensed under the Apache V2 License. See [LICENSE](LICENSE) for more information.
