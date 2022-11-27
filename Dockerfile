FROM golang:1.19-alpine AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /app/echelon
ENTRYPOINT ["/app/echelon", "worker"]
