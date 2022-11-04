package gw

import "cti/erro"

var (
	ErrQueryStringIsRequired = erro.NewError("QUERY_STRING_REQUIRED", "query string is required", nil)
	ErrQueryStringInvalid    = erro.NewError("QUERY_STRING_INVALID", "query string is invalid", nil)
	ErrNoDataSourceAvailable = erro.NewError("NO_DATA_SOURCE_AVAILABLE", "no data source available", nil)
)
