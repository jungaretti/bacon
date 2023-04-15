# Bacon 🥓

Bacon deploys your DNS records from simple config files. You can use Bacon to deploy, backup, and restore your DNS records. Bored already? Let me try again...

Bacon is a peek into the future of DNS record management. It lets you codify your DNS records and deploy them with GitHub Actions. You can use issues, pull requests, and `git blame` to keep track of how and why your DNS records were created. In an emergency, you can use `git revert` to undo recent changes and fix your website.

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

### Authentication

Bacon only works with some DNS providers. Pull requests to add new providers are always welcome!

#### Porkbun

Sign into Porkbun's website and [generate a new API keyset](https://porkbun.com/account/api) for your account. Read the ["Generating API Keys" section of Porkbun's docs](https://kb.porkbun.com/article/190-getting-started-with-the-porkbun-dns-api) for more detailed instructions.

Export two environment variables with your Porkbun API keys:

```bash
export PORKBUN_API_KEY="pk1_123abc789xyz"
export PORKBUN_SECRET_KEY="sk1_xyz123abc789"
```

## Usage

Bacon only offers two subcommands:

- `ping` to double-check your API keys (stored in environment variables)
- `deploy <config>` to deploy DNS records from a YAML config file

### Commands

#### `ping`

Verifies your API keys by pinging Porkbun.

#### `deploy`

Deploys records from a domain's config file by deleting unknown records and creating new records. Add `--delete` to delete outdated records and `--create` to create new records.

### Modes

#### Dry Run

Bacon defaults to its dry-run mode. Execute `bacon deploy` without any flags to preview what it'll do:

```bash
bacon deploy dns/example-com.yml
```

```
Would delete 2 records:
- {225823316 example.com A 123.456.789.112 600 0 }
- {225823318 www.example.com A 123.456.789.112 600 0 }
Would create 2 records:
- { example.com A 789.112.123.456 600 0 }
- { www.example.com A 789.112.123.456 600 0 }
Mock deployment complete
```

#### Modify

Use the `--delete` flag to delete outdated records and the `--create` flag to create new records:

```bash
bacon deploy dns/example-com.yml --delete --create
```

```txt
Deleting 2 records...
- {225823316 example.com A 123.456.789.112 600 0 }
- {225823318 www.example.com A 123.456.789.112 600 0 }
Creating 2 records...
- 225823565
- 225823566
Deployment complete!
```

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

## Built With

- [Cobra](https://cobra.dev/)
- [Porkbun API](https://porkbun.com/api/json/v3/documentation)
