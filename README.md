# Bacon

Tasty DNS record manager for Porkbun

![Demo of Bacon](https://user-images.githubusercontent.com/19893438/166407773-8eb27040-7438-48fc-b5b5-435298b6eb63.gif)

Bacon deploys DNS records from a YAML config file. You can use Bacon to check DNS records into source control and deploy them to Porkbun with GitHub Actions.

## Getting Started

### Installation

Clone this repo and use `make` to build `bin/bacon`.

### Authentication

[Generate a new API keyset online.](https://porkbun.com/account/api)
Next, export two environment variables with your Porkbun API keys:

- `PORKBUN_API_KEY` for your API key
- `PORKBUN_SECRET_KEY` for your secret key

```shell
export PORKBUN_API_KEY="pk1_123abc789xyz"
export PORKBUN_SECRET_KEY="sk1_xyz123abc789"
```

## Usage

Bacon only offers two subcommands:

- `ping` to double-check your API keys (stored in environment variables)
- `deploy <config>` to deploy DNS records from a YAML config file

### `ping`

Verifies your API keys by pinging Porkbun.

### `deploy`

Deploys records from a domain's config file by deleting unknown records and creating new records. Use `-d` to delete existing records and `-c` to create new records.

#### Dry Mode

Bacon defaults to its dry run mode. Call Bacon without any flags to preview what it'll do:

```shell
bacon deploy borkbork.buzz.yml
```

```txt
Would delete 2 records:
- {225823316 borkbork.buzz TXT hello 600 0 }
- {225823318 borkbork.buzz TXT world 600 0 }
Would create 2 records:
- { borkbork.buzz TXT hotdog 600 0 }
- { borkbork.buzz TXT burger 600 0 }
Mock deployment complete
```

#### Deployment Mode

Use the `--delete` flag to delete existing records and the `--create` flag to create new records:

```shell
bacon deploy borkbork.buzz.yml --delete --create
```

```txt
Deleting 2 records...
- {225823316 borkbork.buzz TXT hello 600 0 }
- {225823318 borkbork.buzz TXT world 600 0 }
Creating 2 records...
- 225823565
- 225823566
Partial deployment complete!
```

## Configuration

Store your DNS configuration in a simple YAML file. See [`config.example.yml`](https://github.com/jungaretti/bacon/blob/main/config.example.yml) for a complete example.

```yaml
domain: borkbork.buzz
records:
  - type: CNAME
    host: www.borkbork.buzz
    content: pixie.porkbun.com
    ttl: 600
  - type: MX
    host: borkbork.buzz
    content: fwd2.porkbun.com
    ttl: 600
    priority: 20
  - type: TXT
    host: borkbork.buzz
    content: fizzbuzz
    ttl: 600
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
