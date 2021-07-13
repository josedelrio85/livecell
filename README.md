# Livelead API

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

```bash
curl -G /
-d "lea_id=12345&cat=6&subcat=341&queue=78&ws_id={{wsid}}&ord_id={{order}}&is_client=0&phone=666666666&url={{url}}"
"http://livelead-pre.josedelrio85.me/lead/live"
```

## Build and run dockerfile

You will need a database to use this component, not included yet in dockerfile

```bash
docker image build -t livelead:[version] .
docker container run -d --name livelead -p 4500:4500 livelead:[version]
```

## Helm

### Fake install

```bash
helm install --dry-run --debug --namespace [namespace_name] ./helm-package \
-f ./helm-package/values-pre.yaml | -f ./helm-package/values-pro.yaml
```

### Install

```bash
helm install --name [name] --namespace [namespace_name] ./helm-package \
-f ./helm-package/values-pre.yaml | -f ./helm-package/values-pro.yaml
```

*[pre] => livelead-pre
*[pro] => livelead

### List helm charts

```bash
helm ls --all
```

### Rollback

```bash
helm rollback [helm_name] [revision_number]
```

### Delete

```bash
helm delete [helm_name]
```

```bash
helm del --purge [helm_name]
```
