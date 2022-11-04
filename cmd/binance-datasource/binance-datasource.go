package main

import (
	"cti/ds"
	"log"
	"os"
)

var Version = "-"

func main() {
	log.Printf("version: %s", Version)
	var options []ds.BinanceApiOption
	baseurl, listenAddr := envVars()

	if baseurl != "" {
		options = append(options, ds.BinanceApiBaseUrlOption(baseurl))
	}

	datasource, err := ds.NewBinanceDataSource(options...)
	if err != nil {
		panic(err)
	}

	server := ds.NewDataSourceApiServer(datasource, listenAddr)
	log.Fatalln(server.ListenAndServe())
}

func envVars() (baseurl string, listenAddr string) {
	baseurl = os.Getenv("BINANCE_BASEURL")

	listenAddr = os.Getenv("BINANCE_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":80"
	}
	return
}
