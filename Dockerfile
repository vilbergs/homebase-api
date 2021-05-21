FROM golang:1.16.3-alpine
WORKDIR /go/src/api

RUN apk update && apk add --no-cache gcc musl-dev git bash

# Copy go.mod over before running go mod download will cache dependencies
COPY go.mod go.sum ./
RUN go mod download -x

COPY . .


# RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api \
#     && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate
# CMD ["./bin/api"]

RUN go get github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon --build="go build -o ./bin/api ./cmd/api" --command="./bin/api"