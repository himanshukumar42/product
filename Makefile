proto:
	protoc -I . --go_opt=module=github.com/himanshuk42/product/pkg/pb --go_out=. --go-grpc_opt=module=github.com/himanshuk42/product/pkg/pb --go-grpc_out=. product.proto

server:
	go run cmd/main.go