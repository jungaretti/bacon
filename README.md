# Bacon 🥓

Bacon deploys your DNS records from YAML config files to [Porkbun](https://porkbun.com/api/json/v3/documentation). You can use Bacon to deploy, backup, and restore your DNS records.

![Demo of bacon deploy](https://user-images.githubusercontent.com/19893438/167231076-2f99e0ce-9ed7-40e4-9b1e-fc2fd578cd0f.gif)

## Getting Started

You can use Bacon locally or in a GitHub Actions workflow. For local use, download the latest release from GitHub or build it yourself. For GitHub Actions, use the [Bacon Deploy Action](https://github.com/jungaretti/bacon-deploy-action) in your workflow.

### Installation

[Download the latest release from GitHub.](https://github.com/jungaretti/bacon/releases)

[![Status of release assets](https://github.com/jungaretti/bacon/actions/workflows/release-assets.yml/badge.svg)](https://github.com/jungaretti/bacon/actions/workflows/release-assets.yml)

#### Build it Yourself

1. Install Bacon's prerequisites:
   - [Go](https://go.dev/dl/)
   - [GNU Make](https://ftp.gnu.org/gnu/make/)
2. Clone this repo and use `make` to build `bin/bacon`
3. Authenticate with your DNS provider (see below)

### Authentication

Sign into Porkbun and [generate a new API keyset](https://porkbun.com/account/api) for your account. Read the ["Generating API Keys" section of Porkbun's docs](https://kb.porkbun.com/article/190-getting-started-with-the-porkbun-dns-api) for more detailed instructions. Be sure to enable API access for the domain(s) that you would like to manage with Bacon.

Next, `export` the `PORKBUN_API_KEY` and `PORKBUN_SECRET_KEY` environment variables or add them to an `.env` file. Bacon uses these environment variables to authenticate with Porkbun. If the current directory contains an `.env` file, then Bacon will load its contents into environment variables. See [`.env.example`](https://github.com/jungaretti/bacon/blob/main/.env.example) for an example.

You can use [`bacon ping`](#ping) to check your authentication configuration.

## Usage

Bacon offers a few commands to help you deploy and save your DNS records:

- `ping` to double-check your API keys (stored in environment variables)
- `deploy <config>` to deploy DNS records from a YAML config file
- `print <domain>` to print your DNS records in YAML format

### Commands

#### `ping`

Verifies your API keys by pinging Porkbun.

#### `deploy <config>`

Deploys records from a domain's config file by deleting unknown records and creating new records. Defaults to a dry-run mode that doesn't modify your DNS records.

##### Parameters

- `--delete` disable dry-run deletions and delete outdated records
- `--create` disable dry-run creations and create new records

#### `print <domain>`

Prints records for a domain in YAML format.

##### Notes

Use `>` to redirect output to a Bacon config file. For example, `bacon print example.com > example.com.yml`

## Configuration

See [`config.example.yml`](https://github.com/jungaretti/bacon/blob/main/config.example.yml) for a complete example.

```yaml
domain: example.com
records:
  - type: A
    host: blog.example.com
    content: 123.456.789.112
    ttl: 600
  - type: A
    host: www.example.com
    content: 123.456.789.112
    ttl: 600
    priority: 20
```

### Schema

#### Record

- `type`
- `host`
- `content`
- `ttl`
- `priority`

#### Config

- `domain`
- `records`

## Development

```bash
# Build Bacon
make build

# Run unit tests
make test-unit

# Run system tests (auth required)
make test-system
```

## Built With

- [Cobra](https://cobra.dev/)
- [Porkbun API](https://porkbun.com/api/json/v3/documentation)
