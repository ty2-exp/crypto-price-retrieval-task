package ds

import (
	"net/url"
	_path "path"
	"strings"
)

func UrlParseWithJoin(baseUrl string, path ...string) (*url.URL, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = _path.Join(append([]string{u.Path}, path...)...)

	return u, nil
}

func EscapeDoubleQuote(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}
