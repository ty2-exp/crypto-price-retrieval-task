package ds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BinanceDataSource struct {
	api *BinanceApi
}

func NewBinanceDataSource(options ...BinanceApiOption) (*BinanceDataSource, error) {
	api, err := NewBinanceApi(options...)
	if err != nil {
		return nil, err
	}

	binanceHistorical := &BinanceDataSource{
		api: api,
	}

	return binanceHistorical, nil
}

func (binanceDataSource *BinanceDataSource) Price(symbol string, t time.Time) (float64, error) {
	ts := t.UnixMilli()
	result, err := binanceDataSource.api.Klines(symbol, "1s", ts, 0, 1)
	if err != nil {
		return 0, ErrSourceError.WithAttrs(map[string]any{"err": err})
	}

	if len(result) <= 0 {
		return 0, &ErrNoData
	}

	if len(result[0]) <= 2 {
		return 0, &ErrInvalidResultFormat
	}

	actualTs, ok := result[0][0].(float64)
	if !ok {
		return 0, fmt.Errorf("%w: (expect: float64, actual: %T)",
			ErrResultTypeMismatch.WithAttrs(map[string]any{
				"field":  "ts",
				"expect": "float64",
				"actual": fmt.Sprintf("%T", actualTs), "data": actualTs}),
			actualTs)
	}

	if int64(actualTs) != ts {
		return 0, fmt.Errorf("%w: (expect: %d, actual: %d)",
			ErrResultValueMismatch.WithAttrs(map[string]any{
				"field": "ts", "expect": ts, "actual": actualTs}),
			ts, int64(actualTs))
	}

	openString, ok := result[0][1].(string)
	if !ok {
		return 0, fmt.Errorf("%w: (expect: string, actual: %T)",
			ErrResultTypeMismatch.WithAttrs(map[string]any{
				"field":  "open",
				"expect": "string",
				"actual": fmt.Sprintf("%T", actualTs),
				"data":   openString}),
			openString)
	}

	open, err := strconv.ParseFloat(openString, 64)
	if err != nil {
		return 0, ErrSourceError.WithAttrs(map[string]any{"err": err})
	}

	return open, nil
}

func (binanceDataSource *BinanceDataSource) Average(symbol string, from time.Time, until time.Time, granularity Granularity) (average float64, actualFrom time.Time, actualUntil time.Time, err error) {
	fromTs := from.UnixMilli()
	untilTs := until.UnixMilli()

	if !granularity.IsValid() {
		return 0, time.Time{}, time.Time{}, &ErrInvalidGranularity
	}

	result, err := binanceDataSource.api.Klines(symbol, BinanceApiInterval(granularity), fromTs, untilTs, 1000)
	if err != nil {
		return 0, time.Time{}, time.Time{}, ErrSourceError.WithAttrs(map[string]any{"err": err})
	}

	if len(result) <= 0 {
		return 0, time.Time{}, time.Time{}, &ErrNoData
	}

	sum := 0.0
	for i, dp := range result {
		if len(result[0]) <= 2 {
			return 0, time.Time{}, time.Time{}, ErrInvalidResultFormat.WithAttrs(map[string]any{"arr": i})
		}

		openString, ok := dp[1].(string)
		if !ok {
			return 0, time.Time{}, time.Time{}, fmt.Errorf("%w: (expect: string, actual: %T)",
				ErrResultTypeMismatch.WithAttrs(map[string]any{
					"field":  "open",
					"expect": "string",
					"actual": fmt.Sprintf("%T", openString), "data": openString}),
				openString)
		}

		open, err := strconv.ParseFloat(openString, 64)
		if err != nil {
			return 0, time.Time{}, time.Time{}, ErrDataParseError.WithAttrs(map[string]any{"field": fmt.Sprintf("open[%d]", i), "err": err.Error()})
		}

		sum += open
	}

	actualFromTsFloat, ok := result[0][0].(float64)
	if !ok {
		return 0, time.Time{}, time.Time{}, fmt.Errorf("%w: (expect: float64, actual: %T)",
			ErrResultTypeMismatch.WithAttrs(map[string]any{
				"field":  "from",
				"expect": "float64",
				"actual": fmt.Sprintf("%T", actualFromTsFloat),
				"data":   actualFromTsFloat}),
			actualFromTsFloat)
	}

	actualFromTs := int64(actualFromTsFloat)

	actualUntilTs := int64(0)
	if len(result) >= 2 {
		actualUntilTsFloat, ok := result[len(result)-1][0].(float64)
		if !ok {
			return 0, time.Time{}, time.Time{}, fmt.Errorf("%w: (expect: until, actual: %T)",
				ErrResultTypeMismatch.WithAttrs(map[string]any{
					"field":  "until",
					"expect": "float64",
					"actual": fmt.Sprintf("%T", actualUntilTsFloat),
					"data":   actualUntilTsFloat}),
				actualUntilTsFloat)
		}
		actualUntilTs = int64(actualUntilTsFloat)
	} else {
		actualUntilTs = actualFromTs
	}

	average = sum / float64(len(result))

	return average, time.UnixMilli(actualFromTs), time.UnixMilli(actualUntilTs), nil
}

type BinanceApiOption func(*BinanceApi) error

func BinanceApiBaseUrlOption(baseUrl string) BinanceApiOption {
	return func(binanceApi *BinanceApi) error {
		u, err := url.Parse(baseUrl)
		if err != nil {
			return err
		}

		if u.Scheme != "https" && u.Scheme != "http" {
			return fmt.Errorf("%w: %s", ErrUnsupportedProtocolScheme.WithAttrs(map[string]any{"protocol": u.Scheme}), u.Scheme)
		}

		binanceApi.baseUrl = baseUrl
		return nil
	}
}

func BinanceApiHttpClientOption(httpClient *http.Client) BinanceApiOption {
	return func(binanceApi *BinanceApi) error {
		binanceApi.httpClient = httpClient
		return nil
	}
}

type BinanceApiInterval string

const (
	BinanceApiInterval1s BinanceApiInterval = "1s"
	BinanceApiInterval1m BinanceApiInterval = "1m"
	BinanceApiInterval1h BinanceApiInterval = "1h"
	BinanceApiInterval1d BinanceApiInterval = "1d"
	BinanceApiInterval1M BinanceApiInterval = "1M"
)

func (granularity BinanceApiInterval) isValid() bool {
	switch granularity {
	case BinanceApiInterval1s:
		fallthrough
	case BinanceApiInterval1m:
		fallthrough
	case BinanceApiInterval1h:
		fallthrough
	case BinanceApiInterval1d:
		fallthrough
	case BinanceApiInterval1M:
		return true
	}
	return false
}

type BinanceApi struct {
	baseUrl    string
	httpClient *http.Client
}

func NewBinanceApi(options ...BinanceApiOption) (*BinanceApi, error) {
	binanceApi := &BinanceApi{
		baseUrl: "https://api.binance.us/",
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxConnsPerHost: 10,
				MaxIdleConns:    100,
			},
			Timeout: time.Second * 60,
		},
	}

	for _, option := range options {
		err := option(binanceApi)
		if err != nil {
			return nil, err
		}
	}

	return binanceApi, nil
}

func (api *BinanceApi) Klines(symbol string, interval BinanceApiInterval, startTime int64, endTime int64, limit int) ([][]any, error) {
	u, err := UrlParseWithJoin(api.baseUrl, "api/v3/klines")
	if err != nil {
		return nil, ErrDataParseError.WithAttrs(map[string]any{"field": "url", "err": err})
	}

	if !interval.isValid() {
		return nil, ErrDataParseError.WithAttrs(map[string]any{"field": "interval", "err": "invalid"})
	}

	query := u.Query()
	query.Add("symbol", symbol)
	query.Add("interval", string(interval)) // 1s
	query.Add("startTime", strconv.FormatInt(startTime, 10))
	if endTime > 0 {
		query.Add("endTime", strconv.FormatInt(endTime, 10))
	}
	query.Add("limit", strconv.FormatInt(int64(limit), 10))
	u.RawQuery = query.Encode()

	resp, err := api.httpClient.Get(u.String())
	if err != nil {
		return nil, ErrRequestFailed.WithAttrs(map[string]any{"err": err})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%w: bad status code", ErrBadStatusCode.WithAttrs(map[string]any{"statusCode": resp.StatusCode, "resp": string(body)}))
	}

	var result [][]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, ErrDataParseError.WithAttrs(map[string]any{"field": "result", "err": err})
	}

	return result, nil

}
