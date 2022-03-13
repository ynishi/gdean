# gdean
grcp decide analytics service

## build
### protoc
```
cd pb
curl -o google/rpc/status.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/rpc/status.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative gdean.proto
```
### wire
```
cd cmd/{api,analyzeapi...}
wire
```
### docker
```
docker build -t gdean-analyze -f Dockerfile.analyze .
```
## run server
```
source calc/venv/bin/activate
go run cmd/api/main.go cmd/api/wire_gen.go
```
or
```
docker run -p 50051:50051 gdean-analyze
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
