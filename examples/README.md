# Examples

This directory contains practical examples demonstrating how to use the go-at library for email delivery with templating and styling.

## Running the Example

To run the main example:

```bash
cd examples
go run main.go
```

**Note:** The example uses placeholder SendGrid credentials and will intentionally fail when attempting to send the actual email. This demonstrates the complete workflow without requiring real credentials.

## What the Example Demonstrates

### 1. SendGrid Service Setup
- How to create and configure a SendGrid service
- Setting up the global sender service
- Proper credential management (with placeholders)

### 2. Email Templating
- Creating HTML templates with embedded CSS
- Creating plain text templates
- Using Go template syntax for dynamic content
- Responsive email design with modern CSS

### 3. Template Data Structure
- Defining data structures for template variables
- Organizing data for complex email templates
- Best practices for template data organization

### 4. Email Rendering and Sending
- Rendering templates with actual data
- Handling template rendering errors
- Sending emails using the global service
- Error handling for email delivery

## Template Features

The example includes a professional welcome email template with:

### CSS Features
- **Responsive Design**: Mobile-friendly layout that adapts to different screen sizes
- **Modern Styling**: Clean, professional appearance with gradients and shadows
- **Typography**: Carefully chosen fonts and spacing for readability
- **Interactive Elements**: Styled buttons with hover effects
- **Color Scheme**: Consistent branding with purple gradient theme
- **Layout Components**: Header, content sections, info boxes, and footer

### Template Variables
- `{{.CompanyName}}` - Your company/organization name
- `{{.User.Name}}` - Recipient's full name
- `{{.User.Email}}` - Recipient's email address
- `{{.User.Username}}` - Recipient's username
- `{{.LoginURL}}` - Link to login page
- `{{.SupportURL}}` - Link to support/help page

## Customizing for Your Use Case

### 1. Replace Placeholder Credentials
```go
// Replace these with your actual SendGrid credentials
apiKey := "SG.your-actual-api-key-here"
senderName := "Your Actual Company Name"
senderEmail := "your-actual-email@yourcompany.com"
```

### 2. Customize the Template
- Modify the HTML template to match your brand colors
- Add your company logo
- Adjust the content and messaging
- Add additional template variables as needed

### 3. Extend the Data Structure
```go
type WelcomeData struct {
    User        User
    CompanyName string
    LoginURL    string
    SupportURL  string
    // Add your custom fields here
    ProductName string
    Features    []string
    // etc.
}
```

## Production Considerations

1. **Security**: Store SendGrid API keys in environment variables or secure configuration
2. **Error Handling**: Implement proper error handling and retry logic
3. **Rate Limiting**: Be aware of SendGrid's rate limits
4. **Testing**: Use different templates for development/staging environments
5. **Analytics**: Consider adding tracking pixels or UTM parameters to links
6. **Compliance**: Ensure your emails comply with CAN-SPAM, GDPR, etc.

## Additional Resources

- [SendGrid Documentation](https://docs.sendgrid.com/)
- [Go Template Documentation](https://pkg.go.dev/text/template)
- [Email Design Best Practices](https://sendgrid.com/blog/email-design-best-practices/)
- [CSS in Email Guidelines](https://www.campaignmonitor.com/css/)
