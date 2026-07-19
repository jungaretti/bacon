# Bacon 🥓

Bacon deploys your DNS records from YAML config files to [Porkbun](https://porkbun.com/api/json/v3/documentation). You can use Bacon to deploy, back up, and restore your DNS records.

<img width="1120" height="893" alt="Demo of bacon deploy" src="https://github.com/user-attachments/assets/60ceb3f0-f7fe-4a8b-aac3-57ed1b14b48f" />

## Getting Started

You can use Bacon locally or in a GitHub Actions workflow. For local use, download the latest release from GitHub or build it yourself. For GitHub Actions, use [jungaretti/bacon-deploy-action](https://github.com/jungaretti/bacon-deploy-action) in your workflow.

See [jungaretti/dns](https://github.com/jungaretti/dns) for an example of DNS records managed and deployed with [jungaretti/bacon-deploy-action](https://github.com/jungaretti/bacon-deploy-action).

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

Next, `export` the `PORKBUN_API_KEY` and `PORKBUN_SECRET_KEY` environment variables or add them to a `.env` file. Bacon uses these environment variables to authenticate with Porkbun. If the current directory contains a `.env` file, then Bacon will load its contents into environment variables. See [`.env.example`](https://github.com/jungaretti/bacon/blob/main/.env.example) for an example.

You can use [`bacon ping`](#ping) to check your authentication configuration.

## Usage

Bacon offers a few commands to help you deploy and save your DNS records:

- `ping` to double-check your API keys (stored in environment variables)
- `deploy <config-file>` to deploy DNS records from a YAML config file
- `print <domain>` to print your DNS records in YAML format

### Commands

#### `ping`

Verifies your API keys by pinging Porkbun.

#### `deploy <config-file>`

Deploys records from a config file by deleting, updating, and creating records. Defaults to a dry-run mode that doesn't modify your DNS records.

##### Parameters

- `--dry-run` preview the deployment without making changes
- `--force` execute the deployment without confirmation
- `--output`, `-o` output format: `table` or `json` (default `table`)

#### `print <domain>`

Prints records for a domain in YAML format.

##### Notes

Use `>` to redirect output to a Bacon config file. For example, `bacon print example.com > example.com.yml`.

## Configuration

See [`config.example.yml`](https://github.com/jungaretti/bacon/blob/main/config.example.yml) for a complete example.

```yaml
domain: example.com
records:
  - host: example.com
    type: ALIAS
    ttl: 600
    content: pixie.porkbun.com
  - host: '*.example.com'
    type: CNAME
    ttl: 600
    content: pixie.porkbun.com
  - type: MX
    host: example.com
    content: in1-smtp.messagingengine.com
    ttl: 600
    priority: 10
    notes: Email server
  - type: MX
    host: example.com
    content: in2-smtp.messagingengine.com
    ttl: 600
    priority: 20
    notes: Email server
  - type: A
    host: api.example.com
    content: 42.42.42.42
    ttl: 600
    notes: API endpoint
```

### Schema

#### Record

- `type` - Required. Allowed values: `A`, `MX`, `CNAME`, `ALIAS`, `TXT`, `AAAA`, `SRV`, `TLSA`, `CAA`, `HTTPS`, `SVCB`.
- `host` - Required.
- `content` - Required.
- `ttl` - Required. Minimum value: `600`.
- `priority` - Optional. Allowed for `MX` and `SRV` records.
- `notes` - Optional. Private notes in Porkbun that are not exposed via DNS.

> Bacon ignores `NS` records and records whose host begins with `_acme-challenge`.

#### Config

- `domain` - Required.
- `records` - Optional. DNS records to deploy for the domain.

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
