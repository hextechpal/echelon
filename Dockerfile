FROM golang:1.19-alpine

WORKDIR /app
ADD release/${GOOS}/${GOARCH}/echelon /app/echelon
