server:
	go run cmd/main.go server

migration:
	go run cmd/migrations/main.go migration

##

protoc:
	protoc --go_out=. \
		--go_opt=paths=source_relative \
        --go-grpc_out=. \
        --go-grpc_opt=paths=source_relative \
        pkg/proto/*.proto

build:
	docker build --tag app-github:multistage -f Dockerfile .

run:
	docker run --network=cloud --name app-github --publish 50053:50053 app-github:multistage

start:
	docker start app-github