#!/bin/sh

# this file helps you setup the env vars in the local dev
# copy this file to "./set-local-env.sh", config the env below in "./set-local-env.sh" and run this command in a terminal
# source ./set-local-env.sh

# binance-datasource env vars
export BINANCE_LISTEN_ADDR=:8081

# influxdb-datasource env vars
export IDB_SERVER_URL=https://ap-southeast-2-1.aws.cloud2.influxdata.com
export IDB_ORG=
export IDB_TOKEN=
export IDB_BUCKET=crypto
export IDB_LISTEN_ADDR=:8082

# datasource-gw env vars
export GW_PRICE_DATASOURCE=influxdb:http://127.0.0.1:8082,binance:http://127.0.0.1:8081
export GW_AVERAGE_DATASOURCE=influxdb:http://127.0.0.1:8082,binance:http://127.0.0.1:8081
export GW_SYMBOL=BTCUSD
export GW_LISTEN_ADDR=:8083

# price-periodic-collector env vars
export PPC_DATASOURCE_BASEURL=https://127.0.0.1:8081
export PPC_INFLUX_SERVER_URL=https://ap-southeast-2-1.aws.cloud2.influxdata.com
export PPC_INFLUX_ORG=
export PPC_INFLUX_TOKEN=
export PPC_INFLUX_BUCKET=crypto
export PPC_INFLUX_SYMBOL=BTCUSD

# testing env vars
export TEST_INFLUX_SERVER_URL=https://ap-southeast-2-1.aws.cloud2.influxdata.com
export TEST_INFLUX_ORG=
export TEST_INFLUX_TOKEN=
export TEST_INFLUX_BUCKET=crypto_test