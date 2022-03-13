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
## run server
```
source calc/venv/bin/activate
go run cmd/api/main.go cmd/api/wire_gen.go
```
## test
```
source calc/venv/bin/activate
cd service
go test
```


