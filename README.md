# Bacon 🥓

Bacon deploys your DNS records from simple config files. You can use Bacon to deploy, backup, and restore your DNS records.

![Demo of bacon deploy](https://user-images.githubusercontent.com/19893438/167231076-2f99e0ce-9ed7-40e4-9b1e-fc2fd578cd0f.gif)

## Getting Started

The easiest way to use Bacon is with the [Bacon Deploy Action](https://github.com/jungaretti/bacon-deploy-action). All you need to do is create a new workflow, `uses` the action, and add some parameters. If you want to use Bacon on your own computer, keep reading!

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

Bacon only works with some DNS providers. Pull requests to add new providers are always welcome!

`bacon` reads secrets from the following sources (in order of precedence):

1. `.env` file in the current directory (see `.env.example` for an example)
2. Environment variables

#### Porkbun

Sign into Porkbun's website and [generate a new API keyset](https://porkbun.com/account/api) for your account. Read the ["Generating API Keys" section of Porkbun's docs](https://kb.porkbun.com/article/190-getting-started-with-the-porkbun-dns-api) for more detailed instructions.

##### Required Secrets

- `PORKBUN_API_KEY`
- `PORKBUN_SECRET_KEY`

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

## Contributing

PRs that add new DNS providers or address [open issues](https://github.com/jungaretti/bacon/issues) are always welcome.

### Development

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
