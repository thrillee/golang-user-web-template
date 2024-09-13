---

# Go Web Template

This Go web template is built using **Golang**, with **Gorm** for ORM, **Atlas** for database migrations, and **Fiber** as the web framework. The template includes essential features like user authentication, account management, and OAuth integration (Google, Facebook, Discord). It's designed to be a solid foundation for building scalable web applications.

## Features

- **User Authentication**: Supports login, logout, account verification (via OTP or link), password reset, and profile updates.
- **Account Verification**: Configurable to verify via OTP or link sent to phone/email. The configuration can be done programmatically in the `OTP-dispatch.go` file:
  ```go
  var dispatchFactory = map[OTPDispatchType]dispatchHandler{
      OTP_DISPATCH_SMS: handleSMSDispatch,
  }

  // Example dispatcher handler
  func getDispatcher(dispatchType OTPDispatchType, codeType OTPCodeType) dispatchHandler {
      dispatcher, ok := dispatchFactory[dispatchType]
      if !ok {
          switch {
          case dispatchType == OTP_DISPATCH_EMAIL && codeType == OTP_LONG:
              return handleLongCodeEmailDispatch
          case dispatchType == OTP_DISPATCH_EMAIL && codeType == OTP_SHORT:
              return handleShortCodeEmailDispatch
          }
      }

      return dispatcher
  }
  ```
- **Profile Management**: Update profile information and change profile/display picture.
- **OAuth Support**: Easily integrate with third-party authentication services, including:
  - Google
  - Facebook
  - Discord
- **Email Configuration**: Customizable email templates and sending options.

## Prerequisites

Make sure you have the following installed:

- Go 1.18+
- PostgreSQL (for the database)
- Fiber framework for Go

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-name>
   ```

2. Setup the environment variables in a `.env` file. Example:

   ```bash
   JWT_SECRET_KEY="openssl generate key"
   DATABASE_URL=postgres://root:<user>@<host>:<port>/<db_name>?sslmode=disable
   AUTH_COOKIE_NAME=app_name
   AUTH_COOKIE_DOMAIN=localhost

   VERIFY_LINK=http://localhost:8080/verify?token=<TOKEN>
   SUPPORT_EMAIL=support@thrillee.me
   SUPPORT_PHONE=23481893413

   GOOGLE_REDIRECT_URL=<GOOGLE_REDIRECT_URL>
   GOOGLE_CLIENT_ID=<GOOGLE_CLIENT_ID>
   GOOGLE_CLIENT_SECRET=<GOOGLE_CLIENT_SECRET>

   FACEBOOK_CLIENT_ID=<FACEBOOK_CLIENT_ID>
   FACEBOOK_CLIENT_SECRET=<FACEBOOK_CLIENT_SECRET>
   FACEBOOK_REDIRECT_URL=<FACEBOOK_REDIRECT_URL>

   DISCORD_REDIRECT=<DISCORD_REDIRECT>
   DISCORD_CLIENT_ID=<DISCORD_CLIENT_ID>
   DISCORD_CLIENT=<DISCORD_CLIENT>

   EMAIL_SENDER_EMAIL=<EMAIL_SENDER_EMAIL>
   EMAIL_USERNAME=<EMAIL_USERNAME>
   EMAIL_PASSWORD=<EMAIL_PASSWORD>
   EMAIL_HOST=<EMAIL_HOST>
   EMAIL_PORT=<EMAIL_PORT>
   ```

3. Install the dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   make run/live
   ```

## Database Migration

This project uses Atlas for migrations. To manage migrations:

1. Create new migrations:
   ```bash
   make makemigrations
   ```

2. Apply migrations:
   ```bash
   make migrate
   ```

## Makefile Commands

Here’s an overview of the available `make` commands:

- `make tidy`: Format the code and clean up Go mod files.
- `make audit`: Run code quality checks (vet, static analysis, vulnerability checks).
- `make test`: Run all unit tests.
- `make run/live`: Run the application with live reloading on file changes.
- `make makemigrations`: Create a new database migration.
- `make migrate`: Apply database migrations.
- `make production/deploy`: Build and deploy the application for production.

## OAuth Setup

To enable OAuth for Google, Facebook, and Discord, configure the respective environment variables with your app credentials. You will need to register your app on these platforms to get the client IDs and secrets.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for more details.

---
