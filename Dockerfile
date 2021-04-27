FROM golang:1.16.3-alpine as builder
COPY go.mod go.sum /go/src/github.com/vilbergs/homebase-api/
WORKDIR /go/src/github.com/vilbergs/homebase-api/
RUN go mod download
COPY . /go/src/github.com/vilbergs/homebase-api/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/homebase-api github.com/vilbergs/homebase-api/

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/vilbergs/homebase-api/build/homebase-api /usr/bin/homebase-api
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/homebase-api"]