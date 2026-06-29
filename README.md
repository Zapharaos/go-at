[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/Zapharaos/go-at)](https://pkg.go.dev/mod/github.com/Zapharaos/go-at)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.24.1-61CFDD.svg?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/Zapharaos/go-at)](https://goreportcard.com/report/github.com/Zapharaos/go-at)
![GitHub License](https://img.shields.io/github/license/Zapharaos/go-at)

![GitHub Release](https://img.shields.io/github/v/release/Zapharaos/go-at)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/Zapharaos/go-at/golang.yml)
[![codecov](https://codecov.io/gh/Zapharaos/go-at/graph/badge.svg?token=FB31YKK4EN)](https://codecov.io/gh/Zapharaos/go-at)


# go-at
Go library for streamlining email delivery with templating capabilities.

**go-at** is an email delivery library designed to make sending emails from Go applications effortless. It combines email delivery capabilities with a robust templating system, allowing you to create beautiful, responsive emails with ease.

### Features

- **Easy integration**: Minimal setup with sensible defaults
- **Email delivery**: Send emails via SendGrid or Brevo with an extensible interface for other providers
- **Templating**: Create dynamic content using Go's template syntax
- **Styling**: Add styles through your template definitions
- **Attachments**: Attach files such as PDFs or calendar invites, including inline images

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
	err = goat.Send(goat.NewEmailMessage(
		"user@example.com",
		"Welcome to Acme Corp!",
		content, // plain text
		content, // HTML (same as plain text in this simple example)
    ))
    if err != nil {
        log.Fatal(err)
    }
}
```

For comprehensive usage examples including template variables, pluralization, and fallback behavior, see the [examples package](./examples/).

### Attachments

Attach files by passing their raw bytes — go-at base64-encodes them for the provider.
Use `WithAttachment` for standard files (PDF, `.ics`, etc.) and `WithInlineAttachment`
to embed an image referenced from the HTML body via `cid:<contentID>`:

```go
pdf, _ := os.ReadFile("invoice.pdf")
logo, _ := os.ReadFile("logo.png")

msg := goat.NewEmailMessage(
    "user@example.com",
    "Your invoice",
    "Please find your invoice attached.",            // plain text
    `<p>Invoice attached</p><img src="cid:logo">`,   // HTML referencing the inline image
).
    WithAttachment("invoice.pdf", "application/pdf", pdf).
    WithInlineAttachment("logo.png", "image/png", logo, "logo")

err := goat.Send(msg)
```

> **Note:** Inline embedding and explicit MIME types are fully supported on SendGrid.
> Brevo (brevo-go v1.1.3) infers the MIME type from the filename and has no Content-ID
> field, so an inline attachment is delivered as a regular attachment.

### Using Brevo instead of SendGrid

go-at also ships with a Brevo implementation of the sender interface. Swap the
service creation for `NewBrevoService` and the rest of the API stays the same:

```go
service := goat.NewBrevoService(
    "your-brevo-api-key",
    "Your Company",
    "no-reply@yourcompany.com",
)
restore := goat.SetSenderService(service)
defer restore()
```

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