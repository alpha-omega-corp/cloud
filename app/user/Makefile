server:
	go run cmd/main.go -modfile local.mod

protoc:
	protoc --go_out=. \
		--go_opt=paths=source_relative \
        --go-grpc_out=. \
        --go-grpc_opt=paths=source_relative \
        proto/*.proto

db_init:
	go run cmd/migrations/main.go db init

db_reset:
	go run cmd/migrations/main.go db reset

build:
	docker build --tag app-user:multistage -f Dockerfile .

run:
	docker run --publish 3000:3000 --name app-user app-user:multistage
