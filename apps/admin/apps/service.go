package apps

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bonsus/go-saas/internal/utils/encryption"
	"github.com/jackc/pgx/v5"
)

type service struct {
	repo repository
}

func NewService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(c context.Context, req Request) (app *App, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)
	req.Url = strings.TrimSpace(req.Url)
	req.BackendUrl = strings.TrimSpace(req.BackendUrl)

	if req.Code == "" {
		errorsMap["code"] = append(errorsMap["code"], "Code is required")
	}
	if check := s.repo.CheckAppCode(req.Code); check {
		errorsMap["code"] = append(errorsMap["code"], "code is already exists")
	}
	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "Name is required")
	}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.Create(req)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func (s *service) Update(c context.Context, req Request, id string) (app *App, errorsMap map[string][]string, err error) {
	_, err = s.repo.Read(id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}

	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)
	req.Url = strings.TrimSpace(req.Url)
	req.BackendUrl = strings.TrimSpace(req.BackendUrl)

	if req.Code == "" {
		errorsMap["code"] = append(errorsMap["code"], "Code is required")
	}
	if check := s.repo.CheckAppCode(req.Code, id); check {
		errorsMap["code"] = append(errorsMap["code"], "code is already exists")
	}
	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "Name is required")
	}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
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

func (s *service) UpdateStatus(c context.Context, req Request, id string) (app *App, errorsMap map[string][]string, err error) {
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
func (s *service) UpdateDb(c context.Context, req Request, id string) (app *App, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	_, err = s.repo.Read(id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	req.DbHost, err = encryption.Encrypt(req.DbHost)
	if err != nil {
		errorsMap["db_host"] = append(errorsMap["db_host"], "DB Host is invalid")
	}
	req.DbPort, err = encryption.Encrypt(req.DbPort)
	if err != nil {
		errorsMap["db_port"] = append(errorsMap["db_port"], "DB Port is invalid")
	}
	req.DbUser, err = encryption.Encrypt(req.DbUser)
	if err != nil {
		errorsMap["db_user"] = append(errorsMap["db_user"], "DB User is invalid")
	}
	req.DbPass, err = encryption.Encrypt(req.DbPass)
	if err != nil {
		errorsMap["db_pass"] = append(errorsMap["db_pass"], "DB Pass is invalid")
	}
	req.DbName, err = encryption.Encrypt(req.DbName)
	if err != nil {
		errorsMap["db_name"] = append(errorsMap["db_name"], "DB Name is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.UpdateDb(req, id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func (s *service) Index(c context.Context, param ParamIndex) (result *AppIndex, err error) {
	result, err = s.repo.Index(param)
	return
}

func (s *service) Read(c context.Context, id string) (result *App, err error) {
	result, err = s.repo.Read(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *service) ReadData(c context.Context, id string) (result *AppData, err error) {
	result, err = s.repo.ReadData(id)
	if err != nil {
		return nil, err
	}
	result.DbHost = encryption.DecryptMask(result.DbHost)
	result.DbPort = encryption.DecryptMask(result.DbPort)
	result.DbUser = encryption.DecryptMask(result.DbUser)
	result.DbPass = encryption.DecryptMask(result.DbPass)
	result.DbName = encryption.DecryptMask(result.DbName)
	return result, nil
}

func (s *service) DbTest(c context.Context, id string) (bool, error) {
	result, err := s.repo.ReadData(id)
	if err != nil {
		return false, err
	}
	host, _ := encryption.Decrypt(result.DbHost)
	port, _ := encryption.Decrypt(result.DbPort)
	user, _ := encryption.Decrypt(result.DbUser)
	pass, _ := encryption.Decrypt(result.DbPass)
	db, _ := encryption.Decrypt(result.DbName)
	err = TestPostgresConnection(host, port, user, pass, db)
	if err != nil {
		return false, err
	}
	return true, nil
}
func TestPostgresConnection(host, port, user, password, dbname string) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return fmt.Errorf("connection failed")
	}
	defer conn.Close(ctx)

	return nil
}

func (s *service) Delete(c context.Context, id string) (err error) {
	_, err = s.repo.ReadData(id)
	if err != nil {
		return err
	}
	err = s.repo.Delete(id)
	return
}

func (s *service) PluginIndex(c context.Context, param ParamIndex) (result *AppPluginIndex, err error) {
	result, err = s.repo.PluginIndex(param)
	return
}
func (s *service) PluginCreate(c context.Context, req PluginRequest) (plugin *AppPlugin, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)

	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	// if check := s.repo.CheckPluginCode(req.Code); check {
	// 	errorsMap["code"] = append(errorsMap["code"], "code is already exists")
	// }
	if req.AppId == "" {
		errorsMap["app_id"] = append(errorsMap["app_id"], "AppId is required")
	}
	if check := s.repo.CheckAppId(req.AppId); !check {
		errorsMap["app_id"] = append(errorsMap["app_id"], "AppId is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.PluginCreate(req)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
func (s *service) PluginUpdate(c context.Context, req PluginRequest, Id string) (plugin *AppPlugin, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)

	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	// if check := s.repo.CheckPluginCode(req.Code, Id); check {
	// 	errorsMap["code"] = append(errorsMap["code"], "code is already exists")
	// }

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.PluginUpdate(req, Id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
func (s *service) PluginDelete(c context.Context, id string) (err error) {
	_, err = s.repo.PluginRead(id)
	if err != nil {
		return err
	}
	err = s.repo.PluginDelete(id)
	return
}
func (s *service) PluginRead(c context.Context, id string) (result *AppPlugin, err error) {
	result, err = s.repo.PluginRead(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) ModulIndex(c context.Context, param ParamIndex) (result *AppModulIndex, err error) {
	result, err = s.repo.ModulIndex(param)
	return
}
func (s *service) ModulCreate(c context.Context, req ModulRequest) (result *AppModul, errorsMap map[string][]string, err error) {
	req, errorsMap = validateRequestModul(req)
	if check := s.repo.CheckAppId(req.AppId); !check {
		errorsMap["app_id"] = append(errorsMap["app_id"], "AppId is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err = s.repo.ModulCreate(req)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
func (s *service) ModulUpdate(c context.Context, req ModulRequest, Id string) (result *AppModul, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)

	if req.Code == "" {
		errorsMap["code"] = append(errorsMap["code"], "code is required")
	}
	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	// if check := s.repo.CheckPluginCode(req.Code, Id); check {
	// 	errorsMap["code"] = append(errorsMap["code"], "code is already exists")
	// }
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err = s.repo.ModulUpdate(req, Id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
func (s *service) ModulUpdateStatus(c context.Context, req ModulRequest, Id string) (result *AppModul, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err = s.repo.ModulUpdateStatus(req, Id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
func (s *service) ModulDelete(c context.Context, id string) (err error) {
	_, err = s.repo.ModulRead(id)
	if err != nil {
		return err
	}
	err = s.repo.ModulDelete(id)
	return
}
func (s *service) ModulRead(c context.Context, id string) (result *AppModul, err error) {
	result, err = s.repo.ModulRead(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func validateRequestModul(req ModulRequest) (ModulRequest, map[string][]string) {
	errorsMap := map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)

	if req.Code == "" {
		errorsMap["code"] = append(errorsMap["code"], "code is required")
	}
	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	if req.AppId == "" {
		errorsMap["app_id"] = append(errorsMap["app_id"], "AppId is required")
	}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}

	for i, item := range req.Features {
		prefix := "features." + strconv.Itoa(i)
		if item.Code == "" {
			errorsMap[prefix+".code"] = append(errorsMap[prefix+".code"], "code is required")
		}
		if item.Name == "" {
			errorsMap[prefix+".name"] = append(errorsMap[prefix+".name"], "name is required")
		}
		if item.Status != "active" && item.Status != "inactive" && item.Status != "deleted" {
			errorsMap[prefix+".name"] = append(errorsMap[prefix+".name"], "name is required")
		}
	}
	return req, errorsMap
}

func (s *service) FeatureCreate(c context.Context, req FeatureRequest) (result *AppModulFeature, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)

	if req.Code == "" {
		errorsMap["code"] = append(errorsMap["code"], "code is required")
	}
	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}
	if check := s.repo.CheckAppId(req.AppId); !check {
		errorsMap["app_id"] = append(errorsMap["app_id"], "AppId is invalid")
	}

	if check := s.repo.CheckAppModulId(req.AppModulId); !check {
		errorsMap["app_modul_id"] = append(errorsMap["app_modul_id"], "App Modul Id is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err = s.repo.FeatureCreate(req)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
func (s *service) FeatureUpdate(c context.Context, req FeatureRequest, Id string) (result *AppModulFeature, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)

	if req.Code == "" {
		errorsMap["code"] = append(errorsMap["code"], "code is required")
	}
	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err = s.repo.FeatureUpdate(req, Id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
func (s *service) FeatureUpdateStatus(c context.Context, req FeatureRequest, Id string) (result *AppModulFeature, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err = s.repo.FeatureUpdateStatus(req, Id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
func (s *service) FeatureDelete(c context.Context, id string) (err error) {
	_, err = s.repo.FeatureRead(id)
	if err != nil {
		return err
	}
	err = s.repo.FeatureDelete(id)
	return
}
func (s *service) FeatureRead(c context.Context, id string) (result *AppModulFeature, err error) {
	result, err = s.repo.FeatureRead(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) FeatureBulkDelete(c context.Context, req Request) (success int, failed int, err error) {
	if len(req.Ids) <= 0 {
		return 0, 0, errors.New("ids is required")
	}
	success, failed, err = s.repo.FeatureBulkDelete(req.Ids)
	if err != nil {
		return 0, 0, err
	}

	return success, failed, nil
}
