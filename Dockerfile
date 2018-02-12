### build layer
FROM golang:1.8.1 AS build
WORKDIR /go/src/github.com/turnerlabs/harbor-metrics

# install node
RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -
RUN apt-get install -y nodejs

COPY . .

# install deps
RUN go get -v github.com/golang/dep/cmd/dep
RUN dep ensure -v

# compile server
RUN GOOS=linux GOARCH=386 go build -v -o app .

# build client dashboard widgets
WORKDIR /go/src/github.com/turnerlabs/harbor-metrics/freeboard/widgets/shipmentEnvironmentsByBarge
RUN npm install
RUN npm run build

WORKDIR /go/src/github.com/turnerlabs/harbor-metrics/freeboard/widgets/shipmentEnvironmentsByGroup
RUN npm install
RUN npm run build
###


### run layer
FROM alpine:3.6
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/turnerlabs/harbor-metrics/app .
COPY --from=build /go/src/github.com/turnerlabs/harbor-metrics/freeboard ./freeboard
COPY --from=build /go/src/github.com/turnerlabs/harbor-metrics/freeboard/dashboard.json ./freeboard/dashboard.json
CMD ["./app"]
###