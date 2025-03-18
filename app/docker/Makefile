server:
	go run cmd/main.go -modfile local.mod

build:
	docker build --tag app-docker:multistage -f Dockerfile .

run:
	docker run --publish 3000:3000 --name app-docker app-docker:multistage
