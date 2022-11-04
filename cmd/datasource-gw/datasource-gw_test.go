package main

import (
	"cti/ds"
	"cti/gw"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDataSource(t *testing.T) {
	result := parseDataSources("influxdb:http://127.0.0.1:8082,binance:http://127.0.0.1:8081")

	var expect []ds.DataSourceApiClient
	dsApi, err := ds.NewDefaultDataSourceApiClient("http://127.0.0.1:8082")
	assert.Nil(t, err)
	expect = append(expect, gw.NewDefaultDataSourceApiClient("influxdb", dsApi))
	dsApi, err = ds.NewDefaultDataSourceApiClient("http://127.0.0.1:8081")
	assert.Nil(t, err)
	expect = append(expect, gw.NewDefaultDataSourceApiClient("binance", dsApi))

	assert.Equal(t, expect, result)

	result = parseDataSources("")
	assert.Nil(t, result)
}

func TestParsePriceDataSources(t *testing.T) {
	result := parsePriceDataSources("influxdb:http://127.0.0.1:8082,binance:http://127.0.0.1:8081")

	var expect []ds.PriceDataSourceApi
	dsApi, err := ds.NewDefaultDataSourceApiClient("http://127.0.0.1:8082")
	assert.Nil(t, err)
	expect = append(expect, gw.NewDefaultDataSourceApiClient("influxdb", dsApi))
	dsApi, err = ds.NewDefaultDataSourceApiClient("http://127.0.0.1:8081")
	assert.Nil(t, err)
	expect = append(expect, gw.NewDefaultDataSourceApiClient("binance", dsApi))

	assert.Equal(t, expect, result)

	result = parsePriceDataSources("")
	assert.Nil(t, result)
}

func TestParseAverageDataSources(t *testing.T) {
	result := parseAverageDataSources("influxdb:http://127.0.0.1:8082,binance:http://127.0.0.1:8081")

	var expect []ds.AverageDataSourceApi
	dsApi, err := ds.NewDefaultDataSourceApiClient("http://127.0.0.1:8082")
	assert.Nil(t, err)
	expect = append(expect, gw.NewDefaultDataSourceApiClient("influxdb", dsApi))
	dsApi, err = ds.NewDefaultDataSourceApiClient("http://127.0.0.1:8081")
	assert.Nil(t, err)
	expect = append(expect, gw.NewDefaultDataSourceApiClient("binance", dsApi))

	assert.Equal(t, expect, result)

	result = parseAverageDataSources("")
	assert.Nil(t, result)
}
