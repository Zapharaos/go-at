// Package goat provides email delivery with templating capabilities.
//
// go-at simplifies sending emails from Go applications with built-in templating
// and SendGrid integration.
//
// Basic usage:
//
//	service := goat.NewSendgridService("api-key", "Your Name", "you@company.com")
//	restore := goat.SetSenderService(service)
//	defer restore()
//
//	template := goat.Template{
//		Name:       "welcome",
//		ContentRaw: "Hello {{.Name}}!",
//		Data:       map[string]string{"Name": "John"},
//	}
//	content, _ := template.Render()
//	err := goat.Send("user@example.com", "Welcome", content, content)
//
// For more details, see README.md.
package goat
