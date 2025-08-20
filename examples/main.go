package main

import (
	"fmt"
	"log"

	"github.com/Zapharaos/go-at"
)

// User represents a user in our system
type User struct {
	Name     string
	Email    string
	Username string
}

// WelcomeData contains data for the welcome email template
type WelcomeData struct {
	User        User
	CompanyName string
	LoginURL    string
	SupportURL  string
}

func main() {
	fmt.Println("=== Go-AT Email Library Example ===")
	fmt.Println()

	// Note: This example will fail to send the actual email since we're using
	// placeholder credentials. This is expected and demonstrates the workflow.
	fmt.Println("⚠️  Note: This example uses placeholder credentials and will fail")
	fmt.Println("   to send the actual email. This is expected and demonstrates the workflow.")
	fmt.Println()

	// Step 1: Set up the SendGrid service
	fmt.Println("1. Setting up SendGrid service...")

	// Replace these with your actual SendGrid credentials
	apiKey := "SG.your-api-key-here"
	senderName := "Your Company"
	senderEmail := "no-reply@yourcompany.com"

	service := goat.NewSendgridService(apiKey, senderName, senderEmail)

	// Set the global service
	restore := goat.SetSenderService(service)
	defer restore()

	fmt.Printf("   ✓ SendGrid service configured with sender: %s <%s>\n", senderName, senderEmail)
	fmt.Println()

	// Step 2: Prepare template data
	fmt.Println("2. Preparing template data...")

	user := User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Username: "johndoe",
	}

	data := WelcomeData{
		User:        user,
		CompanyName: "Acme Corp",
		LoginURL:    "https://app.yourcompany.com/login",
		SupportURL:  "https://support.yourcompany.com",
	}

	fmt.Printf("   ✓ Template data prepared for user: %s <%s>\n", user.Name, user.Email)
	fmt.Println()

	// Step 3: Create and render HTML template
	fmt.Println("3. Creating and rendering HTML email template...")

	htmlTemplate := goat.Template{
		Name:       "welcome-html",
		ContentRaw: createWelcomeEmailTemplate(),
		Data:       data,
	}

	htmlContent, err := htmlTemplate.Render()
	if err != nil {
		log.Fatalf("Failed to render HTML template: %v", err)
	}

	fmt.Println("   ✓ HTML template rendered successfully")
	fmt.Printf("   📧 HTML content length: %d characters\n", len(htmlContent))
	fmt.Println()

	// Step 4: Create and render plain text template
	fmt.Println("4. Creating and rendering plain text template...")

	plainTextTemplate := goat.Template{
		Name:       "welcome-text",
		ContentRaw: createPlainTextTemplate(),
		Data:       data,
	}

	plainTextContent, err := plainTextTemplate.Render()
	if err != nil {
		log.Fatalf("Failed to render plain text template: %v", err)
	}

	fmt.Println("   ✓ Plain text template rendered successfully")
	fmt.Printf("   📧 Plain text content length: %d characters\n", len(plainTextContent))
	fmt.Println()

	// Step 5: Send the email
	fmt.Println("5. Attempting to send email...")

	subject := fmt.Sprintf("Welcome to %s, %s!", data.CompanyName, data.User.Name)

	err = goat.Send(user.Email, subject, plainTextContent, htmlContent)
	if err != nil {
		fmt.Printf("   ❌ Expected error occurred (placeholder credentials): %v\n", err)
		fmt.Println("   💡 To actually send emails, replace the placeholder SendGrid credentials")
		fmt.Println("      with your real API key and sender information.")
	} else {
		fmt.Println("   ✅ Email sent successfully!")
	}

	fmt.Println()
	fmt.Println("=== Example completed ===")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("• Replace placeholder SendGrid credentials with real ones")
	fmt.Println("• Customize the email template to match your brand")
	fmt.Println("• Add email validation if needed for your use case")
	fmt.Println("• Integrate into your application workflow")
}

// createWelcomeEmailTemplate returns a complete HTML email template with embedded CSS
func createWelcomeEmailTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome Email</title>
    <style>
        /* Reset and base styles */
        body, table, td, p, h1, h2, h3 {
            margin: 0;
            padding: 0;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;
        }
        
        body {
            background-color: #f8f9fa;
            line-height: 1.6;
            color: #333333;
        }
        
        /* Container */
        .email-container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        }
        
        /* Header */
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 30px 20px;
            text-align: center;
            color: white;
        }
        
        .header h1 {
            font-size: 28px;
            font-weight: 600;
            margin-bottom: 8px;
        }
        
        .header p {
            font-size: 16px;
            opacity: 0.9;
        }
        
        /* Content */
        .content {
            padding: 40px 30px;
        }
        
        .content h2 {
            color: #2c3e50;
            font-size: 24px;
            margin-bottom: 20px;
        }
        
        .content p {
            font-size: 16px;
            margin-bottom: 20px;
            line-height: 1.7;
        }
        
        /* Button */
        .button-container {
            text-align: center;
            margin: 30px 0;
        }
        
        .button {
            display: inline-block;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            text-decoration: none;
            padding: 14px 28px;
            border-radius: 6px;
            font-weight: 600;
            font-size: 16px;
            transition: transform 0.2s ease;
        }
        
        .button:hover {
            transform: translateY(-1px);
        }
        
        /* Info box */
        .info-box {
            background-color: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 20px;
            margin: 25px 0;
            border-radius: 0 6px 6px 0;
        }
        
        .info-box p {
            margin-bottom: 0;
            font-size: 14px;
        }
        
        /* Footer */
        .footer {
            background-color: #f8f9fa;
            padding: 25px 30px;
            text-align: center;
            border-top: 1px solid #e9ecef;
        }
        
        .footer p {
            font-size: 14px;
            color: #6c757d;
            margin-bottom: 10px;
        }
        
        .footer a {
            color: #667eea;
            text-decoration: none;
        }
        
        .footer a:hover {
            text-decoration: underline;
        }
        
        /* Responsive */
        @media only screen and (max-width: 600px) {
            .email-container {
                margin: 0 10px;
            }
            
            .content {
                padding: 30px 20px;
            }
            
            .header {
                padding: 25px 20px;
            }
            
            .header h1 {
                font-size: 24px;
            }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <!-- Header -->
        <div class="header">
            <h1>Welcome to {{.CompanyName}}!</h1>
            <p>We're excited to have you on board</p>
        </div>
        
        <!-- Content -->
        <div class="content">
            <h2>Hello {{.User.Name}}! 👋</h2>
            
            <p>
                Welcome to {{.CompanyName}}! We're thrilled to have you join our community. 
                Your account has been successfully created and you're all set to get started.
            </p>
            
            <p>
                Here are your account details:
            </p>
            
            <div class="info-box">
                <p><strong>Username:</strong> {{.User.Username}}</p>
                <p><strong>Email:</strong> {{.User.Email}}</p>
            </div>
            
            <p>
                Click the button below to log in to your account and start exploring 
                all the features we have to offer.
            </p>
            
            <div class="button-container">
                <a href="{{.LoginURL}}" class="button">Get Started</a>
            </div>
            
            <p>
                If you have any questions or need assistance, don't hesitate to reach out 
                to our support team. We're here to help!
            </p>
        </div>
        
        <!-- Footer -->
        <div class="footer">
            <p>
                Need help? Visit our <a href="{{.SupportURL}}">Support Center</a> 
                or reply to this email.
            </p>
            <p>
                © 2025 {{.CompanyName}}. All rights reserved.
            </p>
        </div>
    </div>
</body>
</html>`
}

// createPlainTextTemplate returns a plain text version of the welcome email
func createPlainTextTemplate() string {
	return `Welcome to {{.CompanyName}}!

Hello {{.User.Name}}!

Welcome to {{.CompanyName}}! We're thrilled to have you join our community. 
Your account has been successfully created and you're all set to get started.

Here are your account details:
- Username: {{.User.Username}}
- Email: {{.User.Email}}

To get started, please visit: {{.LoginURL}}

If you have any questions or need assistance, don't hesitate to reach out 
to our support team at: {{.SupportURL}}

We're here to help!

© 2025 {{.CompanyName}}. All rights reserved.`
}
