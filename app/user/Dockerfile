FROM golang:1.23 AS build-stage

WORKDIR /app

COPY . .
COPY docker.mod ./go.mod
COPY go.sum ./

RUN go get -u ./...
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o /app-user cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app-user /home/nonroot/app-user

EXPOSE 50051

USER nonroot:nonroot

CMD ["/home/nonroot/app-user", "start", "-e", "docker"]