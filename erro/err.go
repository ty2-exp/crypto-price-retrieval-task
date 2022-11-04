package erro

type Error struct {
	error
	Code string         `json:"code"`
	Text string         `json:"err"`
	Attr map[string]any `json:"info"`
}

func NewError(code string, text string, attr map[string]any) Error {
	return Error{Text: text, Code: code, Attr: attr}
}

func (err Error) WithAttrs(attr map[string]any) error {
	e := &err
	e.Attr = attr
	return e
}

func (err Error) Error() string {
	return err.Text
}
