FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app-user cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /migrations-user cmd/migrations/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app-user /home/nonroot/app-user
COPY --from=build-stage /app-user /home/nonroot/migration-user

RUN /app/user/migration-user db init
RUN /app/user/migration-user db reset

EXPOSE 50051

USER nonroot:nonroot

CMD ["/home/nonroot/app-user", "-env", "docker"]