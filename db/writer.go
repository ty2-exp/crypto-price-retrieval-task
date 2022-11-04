package db

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"time"
)

type Writer interface {
	WritePrice(symbol string, price float64, ts time.Time) error
}

type InfluxDbWriter struct {
	client influxdb2.Client
	org    string
	bucket string
}

func NewInfluxDbWriter(serverUrl, org, bucket, token string) *InfluxDbWriter {
	return &InfluxDbWriter{
		client: influxdb2.NewClient(serverUrl, token),
		org:    org,
		bucket: bucket,
	}
}

func (writer *InfluxDbWriter) WritePrice(symbol string, price float64, ts time.Time) error {
	w := writer.client.WriteAPIBlocking(writer.org, writer.bucket)

	p := influxdb2.NewPoint("price",
		map[string]string{"symbol": symbol},
		map[string]interface{}{"open": price},
		ts)

	err := w.WritePoint(context.Background(), p)
	if err != nil {
		return err
	}

	return nil
}
