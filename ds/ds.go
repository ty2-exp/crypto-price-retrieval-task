package ds

import (
	"time"
)

type DataSource interface {
	PriceDataSource
	AverageDataSource
}

type PriceDataSource interface {
	Price(symbol string, ts time.Time) (float64, error)
}

type AverageDataSource interface {
	Average(symbol string, from time.Time, until time.Time, granularity Granularity) (average float64, actualFrom time.Time, actualUntil time.Time, err error)
}

type PriceApiModel struct {
	Price float64 `json:"price"`
}

type PriceAverageApiModel struct {
	Average float64 `json:"average"`
	From    int64   `json:"from"`
	Until   int64   `json:"until"`
}

type Granularity string

const (
	Granularity1s Granularity = "1s"
	Granularity1m Granularity = "1m"
	Granularity1h Granularity = "1h"
	Granularity1d Granularity = "1d"
	Granularity1M Granularity = "1M"
)

func (granularity Granularity) IsValid() bool {
	switch granularity {
	case Granularity1s:
		fallthrough
	case Granularity1m:
		fallthrough
	case Granularity1h:
		fallthrough
	case Granularity1d:
		fallthrough
	case Granularity1M:
		return true
	}

	return false
}

type DataSourceApiClient interface {
	PriceDataSourceApi
	AverageDataSourceApi
}

type PriceDataSourceApi interface {
	Price(symbol string, ts time.Time) (PriceApiModel, error)
}

type AverageDataSourceApi interface {
	Average(symbol string, from time.Time, until time.Time, granularity Granularity) (PriceAverageApiModel, error)
}
