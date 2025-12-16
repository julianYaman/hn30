# build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server

# final stage
FROM alpine:latest

WORKDIR /app

COPY --from=build /server /server

EXPOSE 8080

CMD ["/server"]
