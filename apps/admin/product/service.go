package product

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

type service struct {
	repo repository
}

func NewService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(c context.Context, req Request) (res *Product, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Slug = strings.TrimSpace(req.Slug)

	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}

	if check := s.repo.CheckAppId(req.AppId); !check {
		errorsMap["app_id"] = append(errorsMap["app_id"], "App is invalid")
	}

	// if check := s.repo.CheckAppModulId(req.AppModulId); !check {
	// 	errorsMap["app_modul_id"] = append(errorsMap["app_modul_id"], "App Modul Id is invalid")
	// }

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.Create(req)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func (s *service) Update(c context.Context, req Request, id string) (res *Product, errorsMap map[string][]string, err error) {
	_, err = s.repo.Read(id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Slug = strings.TrimSpace(req.Slug)

	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}
	if check := s.repo.CheckAppId(req.AppId); !check {
		errorsMap["app_id"] = append(errorsMap["app_id"], "AppId is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.Update(req, id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func (s *service) UpdateStatus(c context.Context, req Request, id string) (res *Product, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	_, err = s.repo.Read(id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.UpdateStatus(req, id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func (s *service) Index(c context.Context, param ParamIndex) (result *ProductIndex, err error) {
	result, err = s.repo.Index(param)
	return
}

func (s *service) Read(c context.Context, id string) (result *Product, err error) {
	result, err = s.repo.Read(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (s *service) Delete(c context.Context, id string) (err error) {
	_, err = s.repo.Read(id)
	if err != nil {
		return errors.New("data not found")
	}
	err = s.repo.Delete(id)
	return
}

func (s *service) PriceUpdate(c context.Context, req Request, id string) (res *Product, errorsMap map[string][]string, err error) {
	_, err = s.repo.Read(id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	errorsMap = map[string][]string{}
	validTypes := map[string]bool{
		"trial":        true,
		"monthly":      true,
		"quarterly":    true,
		"semiannually": true,
		"yearly":       true,
		"onetime":      true,
	}
	for i, price := range req.Prices {
		prefix := "prices." + strconv.Itoa(i)
		if !validTypes[price.Type] {
			errorsMap[prefix+".type"] = append(errorsMap[prefix+".type"], "type is invalid")
		}
		if price.Price < 0 {
			errorsMap[prefix+".price"] = append(errorsMap[prefix+".price"], "price must be greater than 0")
		}
		if price.RecurringPrice < 0 {
			errorsMap[prefix+".recurring_price"] = append(errorsMap[prefix+".recurring_price"], "recurring price must be greater than 0")
		}
		if price.Qty <= 0 {
			errorsMap[prefix+".qty"] = append(errorsMap[prefix+".qty"], "qty must be greater than 0")
		}
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.PriceUpdate(req, id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
func (s *service) LimitUpdate(c context.Context, req Request, id string) (res *Product, errorsMap map[string][]string, err error) {
	_, err = s.repo.Read(id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	errorsMap = map[string][]string{}
	for i, limit := range req.Limits {
		prefix := "limits." + strconv.Itoa(i)
		if limit.Name == "" {
			errorsMap[prefix+".name"] = append(errorsMap[prefix+".name"], "name is required")
		}
		if limit.Key == "" {
			errorsMap[prefix+".key"] = append(errorsMap[prefix+".key"], "key is required")
		}
		if limit.Name == "" {
			errorsMap[prefix+".value"] = append(errorsMap[prefix+".value"], "value is required")
		}
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.LimitUpdate(req, id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func (s *service) ModulUpdate(c context.Context, req Request, id string) (res *Product, errorsMap map[string][]string, err error) {
	_, err = s.repo.Read(id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	errorsMap = map[string][]string{}
	for i, item := range req.Moduls {
		prefix := "moduls." + strconv.Itoa(i)
		if check := s.repo.CheckAppModulId(item.AppModulId); !check {
			errorsMap[prefix+".app_modul_id"] = append(errorsMap[prefix+".app_modul_id"], "modul is invalid")
		}
		for n, item2 := range item.Features {
			prefix := "moduls." + strconv.Itoa(i) + ".features." + strconv.Itoa(n)
			// if check := s.repo.CheckProductModulId(item2.ProductModulId); !check {
			// 	errorsMap[prefix+".product_modul_id"] = append(errorsMap[prefix+".product_modul_id"], "modul is invalid")
			// }
			if check := s.repo.CheckAppModulFeatureId(item2.AppModulFeatureId, item.AppModulId); !check {
				errorsMap[prefix+".app_modul_feature_id"] = append(errorsMap[prefix+".app_modul_feature_id"], "feature is invalid")
			}
		}
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.ModulUpdate(req, id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func validateRequestPlan(req Request) (Request, map[string][]string) {
	errorsMap := map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)

	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}

	validTypes := map[string]bool{
		"trial":        true,
		"monthly":      true,
		"quarterly":    true,
		"semiannually": true,
		"yearly":       true,
		"onetime":      true,
	}
	for i, price := range req.Prices {
		prefix := "prices." + strconv.Itoa(i)
		if !validTypes[price.Type] {
			errorsMap[prefix+".type"] = append(errorsMap[prefix+".type"], "type is invalid")
		}
		if price.Price < 0 {
			errorsMap[prefix+".price"] = append(errorsMap[prefix+".price"], "price must be greater than 0")
		}
		if price.RecurringPrice < 0 {
			errorsMap[prefix+".recurring_price"] = append(errorsMap[prefix+".recurring_price"], "recurring price must be greater than 0")
		}
		if price.Qty <= 0 {
			errorsMap[prefix+".qty"] = append(errorsMap[prefix+".qty"], "qty must be greater than 0")
		}
	}
	for i, limit := range req.Limits {
		prefix := "limits." + strconv.Itoa(i)
		if limit.Name == "" {
			errorsMap[prefix+".name"] = append(errorsMap[prefix+".name"], "name is required")
		}
		if limit.Key == "" {
			errorsMap[prefix+".key"] = append(errorsMap[prefix+".key"], "key is required")
		}
		if limit.Name == "" {
			errorsMap[prefix+".value"] = append(errorsMap[prefix+".value"], "value is required")
		}
	}
	return req, errorsMap
}
