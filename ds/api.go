package ds

import (
	"cti/erro"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const DataSourceApiServerRouteGroupV1 = "/api/v1"
const DataSourceApiServerRoutePrice = "/price"
const DataSourceApiServerRouteAverage = "/average"

type ErrorPayload struct {
	Code string         `json:"code"`
	Msg  string         `json:"msg"`
	Attr map[string]any `json:"attr,omitempty"`
}

func NewErrorPayload(err error) ErrorPayload {
	payload := ErrorPayload{Code: "INTERNAL_ERROR", Msg: err.Error()}

	var e *erro.Error
	if errors.As(err, &e) {
		payload.Code = e.Code
		payload.Attr = e.Attr
	}

	return payload
}

func (err ErrorPayload) Error() string {
	return err.Msg
}

type DataSourceApiServer struct {
	listenAddr string
	dataSource DataSource
	Router     *chi.Mux
}

func NewDataSourceApiServer(dataSource DataSource, listenAddr string) *DataSourceApiServer {
	server := &DataSourceApiServer{dataSource: dataSource, listenAddr: listenAddr}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/api/v1", server.v1Route)
	server.Router = r

	return server
}

func (server *DataSourceApiServer) v1Route(r chi.Router) {
	r.Get(DataSourceApiServerRoutePrice, server.Price)
	r.Get(DataSourceApiServerRouteAverage, server.Average)
}

func (server *DataSourceApiServer) Price(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("%w: symbol", ErrDataSourceApiServerQueryStringIsRequired.WithAttrs(map[string]interface{}{"field": "symbol"}))))
		return
	}
	queryTs := r.URL.Query().Get("ts")
	if queryTs == "" {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("%w: ts", ErrDataSourceApiServerQueryStringIsRequired.WithAttrs(map[string]interface{}{"field": "ts"}))))
		return
	}

	ts, err := strconv.ParseInt(queryTs, 10, 64)
	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("convert querystring ts error: %w", err)))
		return
	}

	price, err := server.dataSource.Price(symbol, time.Unix(ts, 0))
	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(err))
		return
	}

	render.Status(r, 200)
	render.JSON(w, r, PriceApiModel{price})
}

func (server *DataSourceApiServer) Average(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("%w: symbol", ErrDataSourceApiServerQueryStringIsRequired.WithAttrs(map[string]any{"field": "symbol"}))))
		return
	}

	queryFrom := r.URL.Query().Get("from")
	if queryFrom == "" {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("%w: from", ErrDataSourceApiServerQueryStringIsRequired.WithAttrs(map[string]any{"field": "from"}))))
		return
	}

	from, err := strconv.ParseInt(queryFrom, 10, 64)
	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("%w: from", ErrDataSourceApiServerQueryStringIsInvalid.WithAttrs(map[string]any{"field": "until", "details": err}))))
		return
	}

	queryUntil := r.URL.Query().Get("until")
	if queryUntil == "" {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("%w: until", ErrDataSourceApiServerQueryStringIsRequired.WithAttrs(map[string]any{"field": "until", "details": err}))))
		return
	}

	until, err := strconv.ParseInt(queryUntil, 10, 64)
	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("%w: until", ErrDataSourceApiServerQueryStringIsInvalid.WithAttrs(map[string]any{"field": "until", "details": err}))))
		return
	}

	granularity := Granularity(r.URL.Query().Get("granularity"))
	if granularity == "" {
		granularity = Granularity1s
	}

	if !granularity.IsValid() {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf(`%w: "%s" granularity is not support`, ErrDataSourceApiServerQueryStringIsInvalid.WithAttrs(map[string]any{"field": "until"}), granularity)))
		return
	}

	result, exactFrom, exactUntil, err := server.dataSource.Average(symbol, time.Unix(from, 0), time.Unix(until, 0), granularity)
	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("request average error: %w", err)))
		return
	}

	render.Status(r, 200)
	render.JSON(w, r, PriceAverageApiModel{result, exactFrom.Unix(), exactUntil.Unix()})
}

func (server *DataSourceApiServer) ListenAndServe() error {
	return http.ListenAndServe(server.listenAddr, server.Router)
}

type DefaultDataSourceApiClient struct {
	baseUrl    string
	httpClient *http.Client
}

func NewDefaultDataSourceApiClient(baseUrl string, options ...DefaultDataSourceApiClientOption) (*DefaultDataSourceApiClient, error) {
	dataSourceApiClient := &DefaultDataSourceApiClient{
		baseUrl: baseUrl,
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxConnsPerHost: 10,
				MaxIdleConns:    100,
			},
			Timeout: time.Second * 30,
		},
	}

	for _, option := range options {
		err := option(dataSourceApiClient)
		if err != nil {
			return nil, err
		}
	}

	return dataSourceApiClient, nil
}

func (client *DefaultDataSourceApiClient) Price(symbol string, ts time.Time) (PriceApiModel, error) {
	u, err := UrlParseWithJoin(client.baseUrl, DataSourceApiServerRouteGroupV1, DataSourceApiServerRoutePrice)
	if err != nil {
		return PriceApiModel{}, err
	}

	query := u.Query()
	query.Add("symbol", symbol)
	query.Add("ts", strconv.Itoa(int(ts.Unix())))
	u.RawQuery = query.Encode()

	resp, err := client.httpClient.Get(u.String())
	if err != nil {
		return PriceApiModel{}, err
	}

	var priceApiModel PriceApiModel
	_, err = client.decodeRespPayload(resp, &priceApiModel)
	if err != nil {
		return PriceApiModel{}, err
	}

	return priceApiModel, nil
}

func (client *DefaultDataSourceApiClient) Average(symbol string, from time.Time, until time.Time, granularity Granularity) (PriceAverageApiModel, error) {
	u, err := UrlParseWithJoin(client.baseUrl, DataSourceApiServerRouteGroupV1, DataSourceApiServerRouteAverage)
	if err != nil {
		return PriceAverageApiModel{}, err
	}

	query := u.Query()
	query.Add("symbol", symbol)
	query.Add("from", strconv.Itoa(int(from.Unix())))
	query.Add("until", strconv.Itoa(int(until.Unix())))
	query.Add("granularity", string(granularity))
	u.RawQuery = query.Encode()

	resp, err := client.httpClient.Get(u.String())
	if err != nil {
		return PriceAverageApiModel{}, err
	}

	var averageApiModel PriceAverageApiModel
	_, err = client.decodeRespPayload(resp, &averageApiModel)
	if err != nil {
		return PriceAverageApiModel{}, err
	}

	return averageApiModel, nil
}

func (client *DefaultDataSourceApiClient) decodeRespPayload(resp *http.Response, model any) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, client.handleErrorPayload(body)
	}

	err = json.Unmarshal(body, &model)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *DefaultDataSourceApiClient) handleErrorPayload(body []byte) ErrorPayload {
	var errPayload ErrorPayload
	err := json.Unmarshal(body, &errPayload)
	if err != nil {
		return NewErrorPayload(fmt.Errorf("unknown error: %s", string(body)))
	}

	return errPayload
}

type DefaultDataSourceApiClientOption func(*DefaultDataSourceApiClient) error

func DefaultDataSourceApiClientHttpClientOption(httpClient *http.Client) DefaultDataSourceApiClientOption {
	return func(binanceApi *DefaultDataSourceApiClient) error {
		binanceApi.httpClient = httpClient
		return nil
	}
}
