package country

import (
	"context"
)

// var validate = validator.New()

type service struct {
	repo repository
}

func NewService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) Countries(c context.Context, params Params) (countries []Country, err error) {
	countries, err = s.repo.Countries(params)
	if err != nil {
		return nil, err
	}
	return countries, nil
}

func (s *service) Provinces(c context.Context, params Params) (provinces []Province, err error) {
	provinces, err = s.repo.Provinces(params)
	if err != nil {
		return nil, err
	}
	return provinces, nil
}

func (s *service) Cities(c context.Context, params Params) (cities *Cities, err error) {
	cities, err = s.repo.Cities(params)
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func (s *service) Districts(c context.Context, params Params) (districts *Districts, err error) {
	districts, err = s.repo.Districts(params)
	if err != nil {
		return nil, err
	}
	return districts, nil
}

func (s *service) Zips(c context.Context, params Params) (zips *Zips, err error) {
	zips, err = s.repo.Zips(params)
	if err != nil {
		return nil, err
	}
	return zips, nil
}

func (s *service) Search(c context.Context, params Params) (searches []Search, err error) {
	searches, err = s.repo.Search(params)
	if err != nil {
		return nil, err
	}
	return searches, nil
}

// func (s *service) Create(c context.Context, req Request) (account *Account, errorsMap map[string][]string, err error) {
// 	errorsMap = map[string][]string{}
// 	err = validate.Struct(req)
// 	if err != nil {
// 		for _, err := range err.(validator.ValidationErrors) {
// 			field := err.Field()
// 			switch field {
// 			case "Name":
// 				errorsMap["name"] = append(errorsMap["no"], "name is required")
// 			case "No":
// 				errorsMap["no"] = append(errorsMap["no"], "no is required")
// 			}
// 		}
// 	}

// 	existingNo, err := s.repo.FindByNo(req.No)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	if existingNo != nil {
// 		errorsMap["no"] = append(errorsMap["no"], "no is exists")
// 	}

// 	if len(errorsMap) > 0 {
// 		return nil, errorsMap, errors.New("")
// 	}

// 	request := FormRequest(req)
// 	result, err := s.repo.Create(request)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return result, nil, nil
// }

// func (s *service) Index(c context.Context, param ParamIndex) (result *AccountIndex, err error) {
// 	param = NewParamIndex(param)

// 	result, err = s.repo.Index(param)
// 	return
// }

// func (s *service) Update(c context.Context, req Request, id string) (account *Account, errorsMap map[string][]string, err error) {
// 	errorsMap = map[string][]string{}
// 	err = validate.Struct(req)
// 	if err != nil {
// 		for _, err := range err.(validator.ValidationErrors) {
// 			field := err.Field()
// 			switch field {
// 			case "Name":
// 				errorsMap["name"] = append(errorsMap["no"], "name is required")
// 			case "No":
// 				errorsMap["no"] = append(errorsMap["no"], "no is required")
// 			}
// 		}
// 	}

// 	existingNo, err := s.repo.FindByNoExceptId(req.No, id)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	if existingNo != nil {
// 		errorsMap["no"] = append(errorsMap["no"], "no is exists")
// 	}

// 	if len(errorsMap) > 0 {
// 		return nil, errorsMap, errors.New("")
// 	}

// 	request := FormRequest(req)
// 	account, err = s.repo.Update(request, id)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return account, nil, nil
// }

// func (s *service) Read(c context.Context, id string) (result *Account, err error) {
// 	result, err = s.repo.Read(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }
