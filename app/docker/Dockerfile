FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app-docker cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app-docker /home/nonroot/app-docker

EXPOSE 3000

USER nonroot:nonroot

CMD ["/home/nonroot/app-docker"]