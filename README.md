<div align="center">
    <img height="128" src="https://github.com/caesar-rocks/docs/raw/master/logo.svg" />
</div>

<div align="center">
    <h1>
        📜 Caesar Mail
    </h1>
</div>

> _Caesar is a Go web framework, designed for productivity. It takes inspiration from traditional web frameworks such as Ruby on Rails, Laravel, Django, Phoenix, AdonisJS, etc._


## Features ✨
- Supports multiple email service providers:
    - SMTP
    - SendGrid
    - Mailgun
    - Amazon SES
    - Resend

- Concurrent email sending using goroutines.
- Configurable options for TLS, SSL, and connection timeouts.
- Keep-alive connection support.
- Graceful shutdown.

## Installation 📦
To use this package in your project, run:

```bash
go get github.com/caesar-rocks/mail
```

## Usage 🚀
Import the Package
```go
import (
    "github.com/caesar-rocks/mail"
)
```

## Configuration 🛠️
Configure your mailer using the `MailConfig` struct:

```go
config := mailer.MailConfig{
    FromName:     "Your Name",
    ReplyToEmail: "reply@example.com",
    Host:         "smtp.example.com",
    HostUser:     "yourusername",
    HostPassword: "yourpassword",
    Port:         "587",
    UseTLS:       true,
    UseSSL:       false,
    Timeout:      10,
    APIService:   mailer.smtp, // Choose the email service provider
    APIKey:       "your-api-key", // Only needed for some services
    KeepAlive:    true,
}

m := mailer.NewMailer(config)
```

## Sending an Email 📤
To send an email, create a MailerMessage and call the Send method:

```go
message := mailer.MailerMessage{
    To:      "recipient@example.com",
    Subject: "Test Subject",
    Body:    "This is a test email.",
}

err := m.Send(message)
if err != nil {
    fmt.Println("Error sending email:", err)
} else {
    fmt.Println("Email sent successfully!")
}
```

## Graceful Shutdown 🔌
To close the mailer and clean up resources, call the Close method:

```go
m.Close()
```

## Configuration Struct 📋
The `MailConfig` struct allows you to configure various settings:

```go
type MailConfig struct {
    FromName     string
    ReplyToEmail string
    Host         string
    HostUser     string
    HostPassword string
    Port         string
    UseTLS       bool
    UseSSL       bool
    Timeout      int
    APIService   APIServiceType
    APIKey       string
    KeepAlive    bool
}
```

| Field          | Type            | Description                                      |
|----------------|-----------------|--------------------------------------------------|
| `FromName`       | string          | The name to be used as the sender.                |
| `ReplyToEmail`   | string          | The reply-to email address.                       |
| `Host`           | string          | The email server host.                            |
| `HostUser`       | string          | The username for the email server.                |
| `HostPassword`  | string          | The password for the email server.                |
| `Port`           | string          | The port of the email server.                     |
| `UseTLS`         | bool            | Boolean indicating whether to use TLS.            |
| `UseSSL`         | bool            | Boolean indicating whether to use SSL.            |
| `Timeout`        | int             | Connection timeout in seconds.                    |
| `APIService`     | APIServiceType  | The email service provider (e.g., smtp, sendgrid, mailgun, amazon-ses, resend). |
| `APIKey`         | string          | The API key for the email service provider (if required). |
| `KeepAlive`      | bool            | Boolean indicating whether to keep the connection alive. |


## Methods 🔧


| Method     | Description                                                         |
|------------|---------------------------------------------------------------------|
| NewMailer  | Creates a new Mail instance with the provided configuration.         |
| Send       | Sends an email using the provided MailerMessage.                     |
| Close      | Closes the mailer and cleans up resources.                           |


## License 📄
This project is licensed under the MIT License. See the [LICENSE](./LICENSE.md) file for details.

Feel free to contribute to this project by opening issues or submitting pull requests. Happy emailing! ✉️