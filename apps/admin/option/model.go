package option

import "encoding/json"

type Option struct {
	Name  string          `json:"name" db:"name"`
	Value json.RawMessage `json:"value" db:"value"`
}

type Request struct {
	Name  string          `json:"name" validate:"required"`
	Value json.RawMessage `json:"value" validate:"required"`
}

func FormRequest(req Request) Request {
	return Request{
		Name:  req.Name,
		Value: req.Value,
	}
}
