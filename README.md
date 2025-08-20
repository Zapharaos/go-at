[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/zapharaos/go-at)](https://pkg.go.dev/mod/github.com/zapharaos/go-at)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.24.1-61CFDD.svg?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/Zapharaos/go-at)](https://goreportcard.com/report/github.com/Zapharaos/go-at)
![GitHub License](https://img.shields.io/github/license/zapharaos/go-at)

![GitHub Release](https://img.shields.io/github/v/release/zapharaos/go-at)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/zapharaos/go-at/golang.yml)
[![codecov](https://codecov.io/gh/Zapharaos/go-at/graph/badge.svg?token=FB31YKK4EN)](https://codecov.io/gh/Zapharaos/go-at)


# go-at
Go library for streamlining email delivery with templating capabilities.

**go-at** is an email delivery library designed to make sending emails from Go applications effortless. It combines email delivery capabilities with a robust templating system, allowing you to create beautiful, responsive emails with ease.

### Features

- **Easy integration**: Minimal setup with sensible defaults
- **Email delivery**: Send emails via SendGrid with an extensible interface for other providers
- **Templating**: Create dynamic content using Go's template syntax
- **Styling**: Add styles through your template definitions

## Installation

```sh
go get github.com/Zapharaos/go-at
```

**Note:** Goat uses [Go Modules](https://go.dev/wiki/Modules) to manage dependencies.

## Usage Examples

Note : This example Sendgrid credential is for demonstration purposes only. Please replace it with your own Sendgrid API key and email address.

```go
package main

import (
    "log"
    "github.com/Zapharaos/go-at"
)

func main() {
    // 1. Set up SendGrid service
    service := goat.NewSendgridService(
        "your-sendgrid-api-key",
        "Your Company",
        "no-reply@yourcompany.com",
    )
    restore := goat.SetSenderService(service)
    defer restore()

    // 2. Create a simple template
    template := goat.Template{
        Name:       "welcome",
        ContentRaw: "Hello {{.Name}}! Welcome to {{.Company}}.",
        Data: map[string]string{
            "Name":    "John Doe",
            "Company": "Acme Corp",
        },
    }

    // 3. Render the template
    content, err := template.Render()
    if err != nil {
        log.Fatal(err)
    }

    // 4. Send the email
    err = goat.Send(
        "user@example.com",
        "Welcome to Acme Corp!",
        content, // plain text
        content, // HTML (same as plain text in this simple example)
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

For comprehensive usage examples including template variables, pluralization, and fallback behavior, see the [examples package](./examples/).

## Development

Install dependencies:
```shell
make dev-deps
```

Run unit tests and generate coverage report:
```shell
make test-unit
```

Run linters:

```shell
make lint
```

Some linter violations can automatically be fixed:

```shell
make fmt
```

## Contributing

We welcome contributions to the go-at library! If you have a bug fix, feature request, or improvement, please open an issue or pull request on GitHub. We appreciate your help in making go-at better for everyone. If you are interested in contributing to the go-at library, please check out our [contributing guidelines](CONTRIBUTING.md) for more information on how to get started.

## License

The project is licensed under the [MIT License](LICENSE).