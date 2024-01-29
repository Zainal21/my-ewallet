# E Wallet

## Getting started

Simple E-Wallet API Service built with [Go Fiber](https://docs.gofiber.io) Golang Framework. (Cooming soon integrated with Payment Gateway Service)

## Dependencies

There is some dependencies that we used in this skeleton:

- [Go Fiber](https://docs.gofiber.io/) [Go Framework]
- [Viper](https://github.com/spf13/viper) [Go Configuration]
- [Cobra](https://github.com/spf13/cobra) [Go Modern CLI]
- [Logrus Logger](https://github.com/sirupsen/logrus) [Go Logger]
- [Goose Migration](https://github.com/pressly/goose) [Go Migration]
- [Gobreaker](https://github.com/sony/gobreaker) [Go Circuit Breaker]

## Features

- Authentication
- Signature JWT
- Get Balance
- Top Up Deposit
- Transfer Deposit
- Get History Transaction
- Integrated with Payment Gateway Service
- Consume API Service With NextJS

## Requirement

- Golang version 1.21 or latest
- Database MySQL

## Usage

### Installation

install required dependencies

```bash
make install
```

### Environment Variable

copy or replace .env.example to .env and change based on your .env

```bash
cp .env.example .env
```

### Run Service

run current service after all dependencies installed

```bash
make start
```

## Database Migration

migration up

```bash
go run main.go db:migrate up
```

migration down

```bash
go run main.go db:migrate down
```

migration reset

```bash
go run main.go db:migrate reset
```

migration reset

```bash
go run main.go db:migrate reset
```

migration redo

```bash
go run main.go db:migrate redo
```

migration status

```bash
go run main.go db:migrate status
```

run seeder

```bash
go run main.go db:seed
```

## API Documentation

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://app.getpostman.com/run-collection/9050639-e2b8fc04-7da9-4b58-9112-f3c42d8189e9?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D9050639-e2b8fc04-7da9-4b58-9112-f3c42d8189e9%26entityType%3Dcollection%26workspaceId%3Ddad9b418-12e0-4f61-8247-1f6a07b6151b)
