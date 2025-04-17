server:
	go run cmd/main.go server

##

build:
	docker build --tag api-gateway:multistage -f Dockerfile .

run:
	docker run --network=cloud --publish 3000:3000 --name api-gateway api-gateway:multistage

start:
	docker start api-gateway


