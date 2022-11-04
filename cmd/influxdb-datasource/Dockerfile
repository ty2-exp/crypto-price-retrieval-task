FROM golang:1.18.3-alpine3.16 AS build

RUN apk add --no-cache make git bash

COPY . /src
WORKDIR /src

# install dependencies
RUN go mod download

# build it
RUN make build-influxdb-datasource

# this results in a single layer image
FROM alpine:3.8
WORKDIR /
COPY --from=build /src/bin/influxdb-datasource /bin/influxdb-datasource
CMD ["/bin/influxdb-datasource"]
