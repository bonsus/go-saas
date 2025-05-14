package promotion

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

func (s *service) Create(c context.Context, req Request) (res *Promotion, errorsMap map[string][]string, err error) {
	// errorsMap = map[string][]string{}
	req, errorsMap = validateRequest(req)
	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.Create(req)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func (s *service) Update(c context.Context, req Request, id string) (res *Promotion, errorsMap map[string][]string, err error) {
	_, err = s.repo.Read(id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	req, errorsMap = validateRequest(req)

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.Update(req, id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func (s *service) UpdateStatus(c context.Context, req Request, id string) (res *Promotion, errorsMap map[string][]string, err error) {
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

func (s *service) Index(c context.Context, param ParamIndex) (result *PromotionIndex, err error) {
	result, err = s.repo.Index(param)
	return
}

func (s *service) Read(c context.Context, id string) (result *Promotion, err error) {
	result, err = s.repo.ReadDetail(id)
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

func validateRequest(req Request) (Request, map[string][]string) {
	errorsMap := map[string][]string{}
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)

	if req.Code == "" {
		errorsMap["code"] = append(errorsMap["code"], "code is required")
	}
	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}

	validTypes := map[string]bool{
		"all":     true,
		"product": true,
		"plan":    true,
		"price":   true,
	}
	for i, item := range req.Items {
		prefix := "items." + strconv.Itoa(i)
		if item.ProductType == "" {
			errorsMap[prefix+".product_type"] = append(errorsMap[prefix+".product_type"], "product type is required")
		}
		if item.Type == "" {
			errorsMap[prefix+".type"] = append(errorsMap[prefix+".type"], "type is required")
		}
		if item.Value <= 0 {
			errorsMap[prefix+".value"] = append(errorsMap[prefix+".value"], "value price must be greater than 0")
		}
		if !validTypes[item.ProductType] {
			errorsMap[prefix+".product"] = append(errorsMap[prefix+".product"], "productxx is invalid")
		} else if item.ProductType != "all" && item.ProductId == "" && item.PriceType == "" {
			errorsMap[prefix+".product"] = append(errorsMap[prefix+".product"], "product is invalid")
		}
	}
	return req, errorsMap
}
