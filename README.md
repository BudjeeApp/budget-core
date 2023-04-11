# budget-core

## Introduction
Hello and welcome to my personal budgeting project! The budget-core application provides a set of APIs for creating and managing
categories, transactions, and payment accounts. These APIs are built using Twirp and RPC. 

Please note that this is an ever-changing application, as I continue to build and improve it.
To stay up-to-date with the latest changes, please refer to the README file.

## Setup

### Prerequisites
* Go 1.20
* Postgres 14.5
* Protobuf
* Dbmate
* Docker

To get started with contributing in `budget-core`, You'll need to run the following steps (assuming you're on MacOS):
1. Install Homebrew https://brew.sh/
2. Install Go `brew install go`
3. Install protobuf `brew install protobuf`

### Building and Running

* Copy the `SAMPLE.env` to `.env`
* Run the `postgres` container (`docker compose up -d postgres`)
* Run `make dev_db`
* Run `make run_server`

## Testing

The `repository` package tests are intended to run against the postgres container, the tests will fail if it is not
running, afterwards running `make test` should execute properly.
