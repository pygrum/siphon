# Siphon

Siphon is a cross-platform malware feed designed to enrich
the threat intelligence process. It pulls the latest identified strains from
verified CTI sources into one portal. 
I personally use it to keep up to date with the latest threats, and get a head-start
on analysing malicious samples!

A database of basic sample information is built up as often as an un-indexed sample 
is found by querying threat intelligence APIs. You can view the most recent samples,
sorted by time, and download them from their source.

## Support

Siphon is designed with only Unix host support in mind, however, it is possible to set it up on Windows
using Git Bash, WSL or similar applications.

## Installation

You can download Siphon source code from the [releases](https://github.com/pygrum/siphon/releases/latest) page, or clone it from this URL.
Then, after entering the containing folder, run `bash scripts/install.sh`. 

### Supported Integrations

| Name          | Setup instructions                                                                                                                                                             |
|---------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| MalwareBazaar | MalwareBazaar integration is used to fetch the latest samples seen in the wild. Create an account at https://bazaar.abuse.ch/ and retrieve api key from your account settings. |

For each integration, add an entry into your configuration file using the
`sources new` command. It should look something like this:

```shell
sources new --name MalwareBazaar --api-key <your-api-key> --endpoint https://mb-api.abuse.ch/api/v1/
```

```yaml
refreshrate: 1
sources:
  - name: MalwareBazaar
    apikey: <your-api-key>
    endpoint: https://mb-api.abuse.ch/api/v1/
```

### Changelog

#### v2.0.0

Siphon has introduced honeypot integration! Agents can now be configured and used on decoy hosts to log
information about and cache samples that are used by attackers in real time. These agents provide the same
interface as other integrations - with the ability to query and download recent samples.

See the [docs](https://github.com/pygrum/siphon/blob/main/docs/DOCS.md) for how to build and configure
agents.
