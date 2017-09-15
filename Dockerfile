### build layer
FROM golang:1.8.1 AS build
WORKDIR /go/src/github.com/turnerlabs/harbor-metrics
COPY . .

# install deps
RUN go get -v github.com/golang/dep/cmd/dep
RUN dep ensure -v

# compile
RUN GOOS=linux GOARCH=386 go build -v -o app .
###


### run layer
FROM alpine:3.6
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/turnerlabs/harbor-metrics/app .
COPY --from=build /go/src/github.com/turnerlabs/harbor-metrics/freeboard ./freeboard
CMD ["./app"]
###