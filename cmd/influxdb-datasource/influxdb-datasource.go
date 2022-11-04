package main

import (
	"cti/ds"
	"log"
	"os"
)

var Version = "-"

func main() {
	log.Printf("version: %s", Version)
	serverUrl, org, token, bucket, listenAddr := envVars()
	datasource, err := ds.NewInfluxDbDataSource(serverUrl, org, bucket, token)
	if err != nil {
		panic(err)
	}

	server := ds.NewDataSourceApiServer(datasource, listenAddr)
	log.Fatalln(server.ListenAndServe())
}

func envVars() (serverUrl string, org string, token string, bucket string, listenAddr string) {
	serverUrl = os.Getenv("IDB_SERVER_URL")
	if serverUrl == "" {
		panic("env IDB_SERVER_URL is required")
	}
	org = os.Getenv("IDB_ORG")
	if org == "" {
		panic("env IDB_ORG is required")
	}
	token = os.Getenv("IDB_TOKEN")
	if token == "" {
		panic("env IDB_TOKEN is required")
	}
	bucket = os.Getenv("IDB_BUCKET")
	if bucket == "" {
		panic("env IDB_BUCKET is required")
	}
	listenAddr = os.Getenv("IDB_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":80"
	}
	return
}
