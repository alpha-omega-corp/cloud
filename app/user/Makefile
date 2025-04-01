server:
	go run cmd/main.go server

migration:
	go run cmd/migrations/main.go db migration

protoc:
	protoc --go_out=. \
		--go_opt=paths=source_relative \
        --go-grpc_out=. \
        --go-grpc_opt=paths=source_relative \
        proto/*.proto

build:
	docker build --tag app-user:multistage -f Dockerfile .

run:
	docker run --publish 50051:50051 --name app-user app-user:multistage
