# OFX Downloader

This Go application automates the process of logging into the website, downloading account statements in OFX format, and parsing the OFX file for transaction details.

## Prerequisites

Before running the application, make sure you have the following installed:

- Go
- Playwright
- Required Go packages (see go.mod file)

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/JacquesBeets/go-ned-ofx.git
   cd fnb-ofx-downloader
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Set up your environment variables:**

   Create a `.env` file in the project root and add the following:

   ```dotenv
    USERN=your_username
    PASSWORD=your_password
    WEBSITE=your_website
    WEBSITE_LOGIN_WAIT=your_website_login_wait
    WEBSITE_LOGOUT_WAIT=your_website_login_wait
   ```

## Usage

Run the application:

    ```bash
    go run main.go
    ```

The application will automate the process of logging in, downloading the OFX file, parsing it, and cleaning up the temporary files.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgments

- [ofxgo](https://github.com/aclindsa/ofxgo) - Go library for parsing OFX files.
- [playwright-go](https://github.com/mxschmitt/playwright-go) - Go bindings for Playwright.
