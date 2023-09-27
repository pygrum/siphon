# Siphon

Siphon is a cross-platform malware feed designed to enrich
the threat intelligence process. It pulls the latest identified strains from
verified CTI sources into one portal. 
I personally use it to keep up to date with the latest threats, and get a head-start
on analysing malicious samples!

## Installation

You can either download Siphon from the [releases](https://github.com/pygurm/siphon/releases/latest)
page, or run the application in a Docker container

### Using docker
#### Dependencies
- Docker
- Make

1. Clone the repository: `git clone https://github.com/pygrum/siphon`
2. Enter the cloned repository and run `make build run`