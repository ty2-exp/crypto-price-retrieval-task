package main

import (
	"cti/db"
	"cti/ds"
	"cti/periodic"
	"log"
	"os"
)

var Version = "-"

func main() {
	log.Printf("version: %s", Version)

	dsBaseUrl, serverUrl, org, token, bucket, symbol := envVars()
	apiClient, err := ds.NewDefaultDataSourceApiClient(dsBaseUrl)
	if err != nil {
		panic(err)
	}

	dbWriter := db.NewInfluxDbWriter(serverUrl, org, bucket, token)
	collector := periodic.NewCollector(apiClient, dbWriter, symbol)
	collector.Start()
}

func envVars() (dsBaseUrl string, serverUrl string, org string, token string, bucket string, symbol string) {
	dsBaseUrl = os.Getenv("PPC_DATASOURCE_BASEURL")
	if dsBaseUrl == "" {
		panic("env PPC_DATASOURCE_BASEURL is required")
	}
	serverUrl = os.Getenv("PPC_INFLUX_SERVER_URL")
	if serverUrl == "" {
		panic("env PPC_INFLUX_SERVER_URL is required")
	}
	org = os.Getenv("PPC_INFLUX_ORG")
	if org == "" {
		panic("env PPC_INFLUX_ORG is required")
	}
	token = os.Getenv("PPC_INFLUX_TOKEN")
	if token == "" {
		panic("env PPC_INFLUX_TOKEN is required")
	}
	bucket = os.Getenv("PPC_INFLUX_BUCKET")
	if bucket == "" {
		panic("env PPC_INFLUX_BUCKET is required")
	}
	symbol = os.Getenv("PPC_SYMBOL")
	if symbol == "" {
		panic("env PPC_SYMBOL is required")
	}
	return
}
