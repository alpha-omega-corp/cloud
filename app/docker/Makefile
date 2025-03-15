server:
	go run cmd/main.go

build:
	docker build --tag api-gateway:multistage -f Dockerfile .

run:
	docker run --publish 3000:3000 --name api-gateway api-gateway:multistage
