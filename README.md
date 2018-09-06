# WARNING!
This project is still in early stages and with plenty of bugs. The config format might change as well as the public APIs which is why the documentation is currently at a minimum. Once all that has been finalized this warning will be removed and this tool/project can be used.

---

![Cliff Logo](docs/assets/logo.png "Cliff Logo")

[![Build Status](https://travis-ci.com/egjiri/cliff.svg?branch=master)](https://travis-ci.com/egjiri/cliff)

## Overview
🏔 **cliff** is a toolset which enables the creation of CLIs through a single YAML config file and no coding is required.

The YAML config file defines all the `commands`, including their `name`, `short` and `long` descriptions, number of `args`, `flags`, flag input types, subcommands, and more. All this gets dispayed in the CLI help output.

Each command could also specify a `run` instruction in the YAML config file which can be any bash command and has access to the `args` and `flags`. For even more control, the `run` instruction can be omitted form the YAML config file and built in code using [Golang](https://golang.org/)

## Example
The following is a simple YAML config file of what a subset of the `docker-compose` CLI would look like if it was using **cliff**

```yaml
name: docker-compose
short: Define and run multi-container applications with Docker.
flags:
  - long: file
    short: f
    type: string
    description: Specify an alternate compose file
    default: docker-compose.yml
    global: true
commands:
  - name: build
    short: Build or rebuild services
    run: echo Build TODO!
  - name: up
    short: Create and start containers
    run: echo Up TODO!
  - name: down
    short: Stop and remove containers, networks, images, and volumes
    run: echo Down TODO!
```

This is the output of running the `docker-compose` command

![Cliff Logo](docs/assets/output.png "Cliff Logo")

## Development Setup
1. Install the Go Programming Language. - [Install Instructions](https://golang.org/doc/install)
2. Install the dep go depenendency manager `go get -u github.com/golang/dep/cmd/dep`
3. Get the vendor packages `dep ensure`
4. Install the go-bindata binary `go get -u github.com/jteeuwen/go-bindata/...`
