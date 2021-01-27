# imscli
![CodeQL](https://github.com/adobe/imscli/workflows/CodeQL/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/adobe/imscli)](https://goreportcard.com/report/github.com/adobe/imscli)

This project is a small CLI tool to interact with the IMS API. The goal of this project
is to provide a small tool that can be used to troubleshoot integrations with IMS.

This project is wrapping adobe/ims-go.

## Installation

Build the CLI or download a prebuilt [release](https://github.com/adobe/imscli/releases).

Example:
```
git clone https://github.com/adobe/imscli

go install
```

## Usage

Once installed, you can start reading the integrated help with the help subcommand.

Examples:

```
imscli help

imscli authorize help

imscli authorize user help
```

The complete documentation of the project is available in the [DOCUMENTATION.md](DOCUMENTATION.md) file.

## Contributing

Contributions are welcomed! Read the [Contributing Guide](CONTRIBUTING.md) for more information.

## Licensing

This project is licensed under the Apache V2 License. See [LICENSE](LICENSE) for more information.
