# Snippetbox

A web application for creating and sharing text snippets. This project is based on Alex Edwards' "Let's Go" book and serves as a practical implementation of Go web development concepts.

## Features

- User authentication (signup, login, logout)
- Create and view text snippets
- Secure session management
- CSRF protection
- RESTful routing

## Technologies

- Go
- JavaScript
- MySQL
- HTML/CSS
- Libraries: httprouter, alice, nosurf

## Setup

1. Clone this repository
2. Make sure Go is installed on your system
3. Set up a MySQL database
4. Configure the application settings
5. Build and run:
   ```
   go build ./cmd/web
   ./web
   ```

## Project Structure

- `cmd/web`: Main application code and handlers
- `ui/html`: HTML templates
- `ui/static`: Static assets (CSS, JavaScript)

## Credits

This project is based on Alex Edwards' "Let's Go" book, which teaches web development with Go.
