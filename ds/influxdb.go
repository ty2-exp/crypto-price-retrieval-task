package ds

import (
	"context"
	"errors"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"time"
)

type InfluxDbDataSource struct {
	token  string
	url    string
	org    string
	bucket string
	client influxdb2.Client
}

func NewInfluxDbDataSource(url string, org string, bucket string, token string) (*InfluxDbDataSource, error) {
	influxDbDataSource := &InfluxDbDataSource{
		token:  token,
		org:    org,
		bucket: bucket,
		client: influxdb2.NewClient(url, token),
	}

	return influxDbDataSource, nil
}

func (influxDbDataSource *InfluxDbDataSource) Price(symbol string, ts time.Time) (float64, error) {
	queryAPI := influxDbDataSource.client.QueryAPI(influxDbDataSource.org)
	tRfc3339 := ts.UTC().Format(time.RFC3339)

	query := fmt.Sprintf(`from(bucket: "%s")
				|> range(start: %s)
				|> filter(fn: (r) => r["_measurement"] == "price")
				|> filter(fn: (r) => r["_field"] == "open")
				|> filter(fn: (r) => r["symbol"] == "%s")
				|> first()
				|> filter(fn: (r) => r["_time"] == %s)
			`, EscapeDoubleQuote(influxDbDataSource.bucket), tRfc3339, symbol, tRfc3339)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		if err.Error() == "invalid: error in building plan while starting program: cannot query an empty range" {
			return 0, &ErrNoData
		}
		return 0, err
	}
	defer result.Close()

	var price *float64 = nil
	var actualTime time.Time
	for result.Next() {
		actualTime = result.Record().Time()
		if v, ok := result.Record().Value().(float64); !ok {
			return 0, errors.New("price type is not valid")
		} else {
			price = &v
		}
		break
	}

	if price == nil {
		return 0, errors.New("price is empty")
	}

	if !actualTime.Equal(ts) {
		return 0, errors.New("timestamp mismatch")
	}

	return *price, nil
}

func (influxDbDataSource *InfluxDbDataSource) Average(symbol string, from time.Time, until time.Time, granularity Granularity) (average float64, actualFrom time.Time, actualUntil time.Time, err error) {
	if !granularity.IsValid() {
		return 0, time.Time{}, time.Time{}, &ErrInvalidGranularity
	}

	// only support granularity 1m
	if granularity != Granularity1m {
		return 0, time.Time{}, time.Time{}, ErrInvalidGranularity.WithAttrs(map[string]any{"details": "only support 1m"})
	}

	// check whether from data point is exists
	fromPrice, err := influxDbDataSource.Price(symbol, from)
	if err != nil {
		return 0, time.Time{}, time.Time{}, err
	}

	// from and until is equal, return from data point directly
	if from.Equal(until) {
		return fromPrice, from, until, nil
	}

	// check whether until data point is exists
	_, err = influxDbDataSource.Price(symbol, until)
	if err != nil {
		return 0, time.Time{}, time.Time{}, err
	}

	queryAPI := influxDbDataSource.client.QueryAPI("")

	tfRfc3339 := from.UTC().Format(time.RFC3339)
	tuRfc3339 := until.UTC().Format(time.RFC3339)

	query := fmt.Sprintf(`from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "price")
				|> filter(fn: (r) => r["_field"] == "open")
				|> filter(fn: (r) => r["symbol"] == "%s")
  				|> mean()
			`, EscapeDoubleQuote(influxDbDataSource.bucket), tfRfc3339, tuRfc3339, symbol)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		if err.Error() == "invalid: error in building plan while starting program: cannot query an empty range" {
			return 0, time.Time{}, time.Time{}, &ErrNoData
		}
		return 0, time.Time{}, time.Time{}, err
	}

	var averagePrice *float64
	actualFrom = time.Time{}
	actualUntil = time.Time{}
	for result.Next() {
		if v, ok := result.Record().ValueByKey("_start").(time.Time); !ok {
			return 0, time.Time{}, time.Time{}, errors.New("from type is not valid")
		} else {
			actualFrom = v
		}

		if v, ok := result.Record().ValueByKey("_stop").(time.Time); !ok {
			return 0, time.Time{}, time.Time{}, errors.New("until type is not valid")
		} else {
			actualUntil = v
		}

		if v, ok := result.Record().Value().(float64); !ok {
			return 0, time.Time{}, time.Time{}, errors.New("price type is not valid")
		} else {
			averagePrice = &v
		}
		break
	}

	if averagePrice == nil {
		return 0, actualFrom, actualUntil, errors.New("average is empty")
	}

	if !actualFrom.Equal(from) {
		return 0, actualFrom, actualUntil, errors.New("from timestamp mismatch")
	}

	if !actualUntil.Equal(until) {
		return 0, actualFrom, actualUntil, errors.New("to timestamp mismatch")
	}

	return *averagePrice, actualFrom, actualUntil, nil
}
