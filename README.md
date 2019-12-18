# Noname API

This component is responsible of receiving updates from Leontel environment. When a consultant sets a new categorization in a lead, a GET request is sended to this endpoint and after needed process (decode data, some data completion) stored in a specific environment.

## Why do we need it

We need this component to get the possibility of analyze the data managed in the Smart Center on real time.

## How to run the component

This component has been developed using the following Go! version:

```bash
go version go1.13.4 windows/amd64
```

It is a HTTP service that could be run locally on 4500 port using:

```bash
go run cmd/main.go
```

## How to run the tests

```bash
go test ./...
```

## Example GET request

```bash
curl -G /
-d "phone=666666666&wsid=1234&queue=244798797&lea_id=12345" /
"http://localhost:4500/status/test"
```
