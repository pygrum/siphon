# Siphon

Siphon is a cross-platform malware feed designed to enrich
the threat intelligence process. It pulls the latest identified strains from
verified CTI sources into one portal. 
I personally use it to keep up to date with the latest threats, and get a head-start
on analysing malicious samples!

A database of basic sample information is built up as often as an un-indexed sample 
is found by querying threat intelligence APIs. You can view the most recent samples,
sorted by time, and download them from their source.

## Installation

You can either download Siphon from the [releases](https://github.com/pygurm/siphon/releases/latest)
page, or run the application in a Docker container

### Using docker
#### Dependencies
- Docker
- Make

1. Clone the repository: `git clone https://github.com/pygrum/siphon`
2. Enter the cloned repository and run `make build run`

## Supported Integrations

| Name          | Setup instructions                                                                                                                                                             |
|---------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| MalwareBazaar | MalwareBazaar integration is used to fetch the latest samples seen in the wild. Create an account at https://bazaar.abuse.ch/ and retrieve api key from your account settings. |

For each integration, add an entry into your configuration file. It should look something like
this:

```yaml
refreshrate: 1
sources:
  - name: MalwareBazaar
    apikey: <your-api-key>
    endpoint: https://mb-api.abuse.ch/api/v1/
```