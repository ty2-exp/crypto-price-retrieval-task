package main

import (
	"cti/ds"
	"cti/gw"
	"log"
	"os"
	"strings"
)

var Version = "-"

func main() {
	log.Printf("version: %s", Version)

	envPriceDataSource, envAverageDataSources, listenAddr, symbol := envVars()
	priceDataSources := parsePriceDataSources(envPriceDataSource)
	averageDataSources := parseAverageDataSources(envAverageDataSources)

	server := gw.NewDataSourceApiGw(priceDataSources, averageDataSources, symbol, listenAddr)
	log.Fatalln(server.ListenAndServe())
}

func envVars() (priceDataSource string, averageDataSources string, listenAddr string, symbol string) {
	priceDataSource = os.Getenv("GW_PRICE_DATASOURCE")
	averageDataSources = os.Getenv("GW_AVERAGE_DATASOURCE")
	listenAddr = os.Getenv("GW_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":80"
	}

	symbol = os.Getenv("GW_SYMBOL")
	if symbol == "" {
		panic("env GW_SYMBOL is required")
	}
	return
}

func parseDataSources(str string) []ds.DataSourceApiClient {
	if str == "" {
		return nil
	}

	var dsArr []ds.DataSourceApiClient
	dsInfo := strings.Split(str, ",")

	for _, info := range dsInfo {
		v := strings.SplitN(info, ":", 2)
		dsApi, err := ds.NewDefaultDataSourceApiClient(v[1])
		if err != nil {
			panic(err)
		}

		dsArr = append(dsArr, gw.NewDefaultDataSourceApiClient(v[0], dsApi))
	}

	return dsArr
}

func parsePriceDataSources(str string) []ds.PriceDataSourceApi {
	if str == "" {
		return nil
	}

	var api []ds.PriceDataSourceApi
	dataSources := parseDataSources(str)
	for _, datasource := range dataSources {
		api = append(api, datasource)
	}

	return api
}

func parseAverageDataSources(str string) []ds.AverageDataSourceApi {
	if str == "" {
		return nil
	}

	var api []ds.AverageDataSourceApi
	dataSources := parseDataSources(str)
	for _, datasource := range dataSources {
		api = append(api, datasource)
	}

	return api
}
