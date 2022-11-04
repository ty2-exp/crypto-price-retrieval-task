package ds

import "cti/erro"

var (
	ErrNoData                                   = erro.NewError("NO_DATA", "no data", nil)
	ErrInvalidResultFormat                      = erro.NewError("INVALID_RESULT_FORMAT", "invalid result format", nil)
	ErrDataSourceApiServerQueryStringIsRequired = erro.NewError("QUERY_STRING_REQUIRED", "query string is required", nil)
	ErrDataSourceApiServerQueryStringIsInvalid  = erro.NewError("QUERY_STRING_INVALID", "query string is invalid", nil)
	ErrUnsupportedProtocolScheme                = erro.NewError("UNSUPPORTED_PROTOCOL_SCHEME", "unsupported protocol scheme", nil)
	ErrResultValueMismatch                      = erro.NewError("RESULT_VALUE_MISMATCH", "result value mismatch", nil)
	ErrResultTypeMismatch                       = erro.NewError("RESULT_TYPE_MISMATCH", "result type mismatch", nil)
	ErrInvalidGranularity                       = erro.NewError("INVALID_GRANULARITY", "invalid granularity", nil)
	ErrSourceError                              = erro.NewError("SOURCE_ERROR", "underlying data source error", nil)
	ErrDataParseError                           = erro.NewError("DATA_PARSE_ERROR", "data parse error", nil)
	ErrBadStatusCode                            = erro.NewError("BAD_STATUS_CODE", "bad status code", nil)
	ErrRequestFailed                            = erro.NewError("REQUEST_FAILED", "failed to send request", nil)
)
