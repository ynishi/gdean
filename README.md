# gdean

grcp decide analytics service

## build

### mage

```
mage build
```

### manual

#### protoc

```
cd pb
curl -o google/rpc/status.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/rpc/status.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative gdean.proto
```

#### wire

```
cd cmd/{api,analyzeapi...}
wire
```

#### gqlgen

```
cd gql
gqlgen
```

#### docker

```
docker build -t gdean-analyze -f Dockerfile.analyze .
docker build -t gdean-issue -f Dockerfile.issue .
docker build -t gdean-gql -f Dockerfile.gql .
```

## run server

```
source calc/venv/bin/activate
go run cmd/xxxapi/main.go cmd/xxxapi/wire_gen.go
```

```
cd gql
go run ./server.go
```

or

```
docker run -p 50051:50051 -p 8081:8080 gdean-analyze
docker run -p 50052:50051 -p 8082:8080 gdean-issue
docker run -p 8080:8080 gdean-gql
```

## test

```
source calc/venv/bin/activate
cd service
go test
```

## lint

```
golangci-lint run
```
