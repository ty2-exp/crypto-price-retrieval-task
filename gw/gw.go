package gw

import (
	"cti/ds"
	"cti/erro"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	"time"
)

type DefaultPayload struct {
	Data   any     `json:"data"`
	Source *string `json:"source"`
}

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

type DataSourceApiGw struct {
	priceDataSource   []ds.PriceDataSourceApi
	averageDataSource []ds.AverageDataSourceApi
	router            *chi.Mux
	symbol            string
	listenAddr        string
}

func NewDataSourceApiGw(priceDataSource []ds.PriceDataSourceApi, averageDataSource []ds.AverageDataSourceApi, symbol string, listenAddr string) *DataSourceApiGw {
	server := &DataSourceApiGw{
		priceDataSource:   priceDataSource,
		averageDataSource: averageDataSource,
		symbol:            symbol,
		listenAddr:        listenAddr,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/api/v1", server.v1Route)
	server.router = r

	return server
}

func (server *DataSourceApiGw) v1Route(r chi.Router) {
	r.Get("/price", server.price)
	r.Get("/average", server.average)
}

func (server *DataSourceApiGw) price(w http.ResponseWriter, r *http.Request) {
	queryTs := r.URL.Query().Get("ts")
	if queryTs == "" {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("%w: ts", ErrQueryStringIsRequired.WithAttrs(map[string]interface{}{"field": "ts"}))))
		return
	}

	ts, err := strconv.ParseInt(queryTs, 10, 64)
	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, ErrQueryStringInvalid.WithAttrs(map[string]any{"field": "ts"}))
		return
	}

	errs := make(map[string]error)
	for i, dataSource := range server.priceDataSource {
		var sourceId *string
		if d, ok := dataSource.(DataSourceApiClient); ok {
			s := d.Id()
			sourceId = &s
		}

		result, err := dataSource.Price(server.symbol, time.Unix(ts, 0))
		if err != nil {
			if sourceId != nil {
				errs[*sourceId] = err
			} else {
				errs[fmt.Sprintf("%d", i)] = err
			}
			continue
		}

		render.Status(r, 200)
		render.JSON(w, r, DefaultPayload{result, sourceId})
		return
	}

	render.Status(r, 400)
	render.JSON(w, r, NewErrorPayload(ErrNoDataSourceAvailable.WithAttrs(map[string]any{"errs": errs})))
}

func (server *DataSourceApiGw) average(w http.ResponseWriter, r *http.Request) {
	queryFrom := r.URL.Query().Get("from")
	if queryFrom == "" {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("%w: from", ErrQueryStringIsRequired.WithAttrs(map[string]interface{}{"field": "from"}))))
		return
	}

	from, err := strconv.ParseInt(queryFrom, 10, 64)
	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, ErrQueryStringInvalid.WithAttrs(map[string]any{"field": "from"}))
		return
	}

	queryUntil := r.URL.Query().Get("until")
	if queryUntil == "" {
		render.Status(r, 400)
		render.JSON(w, r, NewErrorPayload(fmt.Errorf("%w: until", ErrQueryStringIsRequired.WithAttrs(map[string]interface{}{"field": "until"}))))
		return
	}

	until, err := strconv.ParseInt(queryUntil, 10, 64)
	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, ErrQueryStringInvalid.WithAttrs(map[string]any{"field": "until"}))
		return
	}

	granularity := ds.Granularity(r.URL.Query().Get("granularity"))
	if granularity == "" {
		granularity = ds.Granularity1s
	}

	errs := make(map[string]error)
	for i, dataSource := range server.averageDataSource {
		var sourceId *string
		if d, ok := dataSource.(DataSourceApiClient); ok {
			s := d.Id()
			sourceId = &s
		}

		result, err := dataSource.Average(server.symbol, time.Unix(from, 0), time.Unix(until, 0), granularity)
		if err != nil {
			if sourceId != nil {
				errs[*sourceId] = err
			} else {
				errs[fmt.Sprintf("%d", i)] = err
			}
			continue
		}

		render.Status(r, 200)
		render.JSON(w, r, DefaultPayload{result, sourceId})
		return
	}

	render.Status(r, 400)
	render.JSON(w, r, NewErrorPayload(ErrNoDataSourceAvailable.WithAttrs(map[string]any{"errs": errs})))
}

func (server *DataSourceApiGw) ListenAndServe() error {
	return http.ListenAndServe(server.listenAddr, server.router)
}
