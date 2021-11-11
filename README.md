# Home Loans Project

This project is about home loans application of a bank so that customers can apply for loans online, then officers can check the submission file and accept or reject the application.

## Prerequisite

To run this program, you will need

### App Dependencies

```$xslt
- Golang 1.12+
- Go mod enabled
```

## How to Run

### Setup App Config

```
cp .env.example .env
```

### Run Application

```
make run
```

## How to Test

```
make test
```

## How to Lint

```
make lint
```

## Deployment

### Build

```
make build
```

## Configuration

| NAME | DESCRIPTION | TYPE | VALUE
| ------ | ------ | ------ | ------ |
| APP_NAME | Application name | string | alphabet |
| APP_PORT | Application port | int | number |
| LOG_LEVEL | Mode for log level configuration | string | debug/info |
| ENVIRONMENT | Application environment | string | development |
| JWT_SECRET | JWT Secret | string | alphabet |
