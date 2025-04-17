FROM golang:1.23 AS build-stage

WORKDIR /app

COPY . .
COPY docker.mod ./go.mod
COPY go.sum ./

RUN go get -u ./...
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /api-gateway /home/nonroot/api-gateway

EXPOSE 3000

USER nonroot:nonroot

CMD ["/home/nonroot/api-gateway", "server", "-e", "docker"]