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
        proto/*.proto


build:
	docker build --tag app-docker:multistage -f Dockerfile .

run:
	docker run --network=cloud --publish 3000:3000 --name app-docker app-docker:multistage

start:
	docker start app-docker