# Bacon

Bacon is a tasty DNS manager for Porkbun.

## Getting Started

Bacon uses Porkbun's API to create, delete, and deploy DNS records.

### Installation

Use `make` to build `bin/bacon`.

## Usage

Export the `PORKBUN_API_KEY` and `PORKBUN_SECRET_KEY` environment variables to authenticate with Porkbun. Use `bacon ping` to test your configuration.

## Built With

- [Cobra](https://cobra.dev/)
- [Porkbun API](https://porkbun.com/api/json/v3/documentation)
