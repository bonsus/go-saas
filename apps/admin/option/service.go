package option

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type service struct {
	repo repository
}

func NewService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) Save(ctx context.Context, req Request) (option *Option, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	err = validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			switch field {
			case "Name":
				errorsMap["name"] = append(errorsMap["name"], "nama wajib diisi")
			case "Value":
				errorsMap["value"] = append(errorsMap["value"], "value wajib diisi")
			}
		}
	}
	if len(errorsMap) > 0 {
		return option, errorsMap, errors.New("")
	}

	formData := FormRequest(req)
	data, err := s.repo.Save(Option(formData))
	if err != nil {
		return option, nil, err
	}

	return data, nil, nil
}
func (s *service) Get(name string) (option Option, err error) {
	data, err := s.repo.Get(name)
	if err != nil {
		return option, err
	}
	option.Value = data
	return option, nil
}
