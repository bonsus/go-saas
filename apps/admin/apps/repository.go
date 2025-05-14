package apps

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// Create(account Request) (*Request, error)
	// Index(param ParamIndex) (*AccountIndex, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

func (r *repository) Create(req Request) (*App, error) {
	Id := uuid.NewString()
	timeNow := time.Now()
	data := App{
		Id:          Id,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		Url:         req.Url,
		BackendUrl:  req.BackendUrl,
		Status:      req.Status,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	result := r.db.Create(&data)

	if result.Error != nil {
		return nil, result.Error
	}
	res, _ := r.Read(Id)
	return res, nil
}

func (r *repository) Update(req Request, Id string) (res *App, err error) {
	updateData := map[string]interface{}{
		"Code":        req.Code,
		"name":        req.Name,
		"description": req.Description,
		"image":       req.Image,
		"url":         req.Url,
		"backend_url": req.BackendUrl,
		// "db_host":     req.DbHost,
		// "db_port":     req.DbPort,
		// "db_user":     req.DbUser,
		// "db_pass":     req.DbPass,
		// "db_name":     req.DbName,
		"status":     req.Status,
		"updated_at": time.Now(),
	}
	result := r.db.Model(&App{}).Where("id = ?", Id).Updates(updateData)
	if result.Error != nil {
		return nil, errors.New("update failed")
	}
	res, _ = r.Read(Id)
	return res, nil
}

func (r *repository) UpdateStatus(req Request, Id string) (res *App, err error) {
	result := r.db.Model(&App{}).Where("id = ?", Id).Updates(map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return nil, result.Error
	}
	res, _ = r.Read(Id)
	return res, nil
}
func (r *repository) UpdateDb(req Request, Id string) (res *App, err error) {
	updateData := map[string]interface{}{
		"db_host":    req.DbHost,
		"db_port":    req.DbPort,
		"db_user":    req.DbUser,
		"db_pass":    req.DbPass,
		"db_name":    req.DbName,
		"updated_at": time.Now(),
	}
	result := r.db.Model(&App{}).Where("id = ?", Id).Updates(updateData)
	if result.Error != nil {
		return nil, errors.New("update failed")
	}
	res, _ = r.Read(Id)
	return res, nil
}

func (r *repository) Index(param ParamIndex) (*AppIndex, error) {
	var data []App
	var total int64

	if param.Page < 1 {
		param.Page = 1
	}
	if param.Perpage < 1 {
		param.Perpage = 20
	}

	offset := (param.Page - 1) * param.Perpage
	dbQuery := r.db.Model(&App{}).Preload("Plugins").Order("created_at DESC")
	if param.Search != "" {
		search := "%" + param.Search + "%"
		dbQuery = dbQuery.Where(
			"name ILIKE ? OR code ILIKE ?",
			search, search,
		)
	}
	dbQuery.Count(&total)
	dbQuery.Limit(param.Perpage).Offset(offset).Find(&data)

	totalPage := int((total + int64(param.Perpage) - 1) / int64(param.Perpage))

	return &AppIndex{
		Data:      data,
		Total:     total,
		TotalPage: totalPage,
		Page:      param.Page,
		PerPage:   param.Perpage,
	}, nil
}

func (r *repository) Read(id string) (*App, error) {
	var app App
	check := r.db.Model(App{}).Where("id = ?", id).Preload("Plugins").Find(&app)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("id not found")
	}
	app.DbHost = "*****"
	app.DbPort = "*****"
	app.DbUser = "*****"
	app.DbPass = "*****"
	app.DbName = "*****"
	return &app, nil
}
func (r *repository) ReadData(id string) (*AppData, error) {
	var app AppData
	check := r.db.Model(App{}).Where("id = ?", id).Find(&app)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("id not found")
	}
	return &app, nil
}

func (r *repository) Delete(id string) error {
	var app App
	if err := r.db.Where("id = ?", id).Find(&app).Error; err != nil {
		return errors.New("id not found")
	}

	delete := r.db.Where("id = ?", id).Delete(&App{})
	if delete.Error != nil {
		return errors.New("data cannot be deleted")
	}

	return nil
}

func (r *repository) PluginIndex(param ParamIndex) (*AppPluginIndex, error) {
	var data []AppPlugin
	var total int64

	if param.Page < 1 {
		param.Page = 1
	}
	if param.Perpage < 1 {
		param.Perpage = 20
	}

	offset := (param.Page - 1) * param.Perpage
	dbQuery := r.db.Model(&AppPlugin{}).Order("created_at DESC")
	if param.Search != "" {
		search := "%" + param.Search + "%"
		dbQuery = dbQuery.Where(
			"name ILIKE ? OR code ILIKE ?",
			search, search,
		)
	}
	dbQuery.Count(&total)
	dbQuery.Limit(param.Perpage).Offset(offset).Find(&data)

	totalPage := int((total + int64(param.Perpage) - 1) / int64(param.Perpage))

	return &AppPluginIndex{
		Data:      data,
		Total:     total,
		TotalPage: totalPage,
		Page:      param.Page,
		PerPage:   param.Perpage,
	}, nil
}
func (r *repository) PluginCreate(req PluginRequest) (*AppPlugin, error) {
	Id := uuid.NewString()
	timeNow := time.Now()
	plugin := AppPlugin{
		Id:        Id,
		AppId:     req.AppId,
		Code:      req.Code,
		Name:      req.Name,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	result := r.db.Create(&plugin)

	if result.Error != nil {
		return nil, result.Error
	}
	res, _ := r.PluginRead(Id)
	return res, nil
}
func (r *repository) PluginUpdate(req PluginRequest, Id string) (*AppPlugin, error) {
	updateData := map[string]interface{}{
		"code":       req.Code,
		"name":       req.Name,
		"updated_at": time.Now(),
	}
	result := r.db.Model(&AppPlugin{}).Where("id = ?", Id).Updates(updateData)
	if result.Error != nil {
		return nil, errors.New("update failed")
	}
	res, _ := r.PluginRead(Id)
	return res, nil
}
func (r *repository) PluginDelete(id string) error {
	if err := r.db.Where("id = ?", id).Find(&AppPlugin{}).Error; err != nil {
		return errors.New("id not found")
	}

	delete := r.db.Where("id = ?", id).Delete(&AppPlugin{})
	if delete.Error != nil {
		return errors.New("data cannot be deleted")
	}

	return nil
}
func (r *repository) PluginRead(id string) (*AppPlugin, error) {
	var plugin AppPlugin
	check := r.db.Model(AppPlugin{}).Where("id = ?", id).Find(&plugin)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("id not found")
	}
	return &plugin, nil
}

func (r *repository) ModulIndex(param ParamIndex) (*AppModulIndex, error) {
	var data []AppModul
	var total int64

	if param.Page < 1 {
		param.Page = 1
	}
	if param.Perpage < 1 {
		param.Perpage = 20
	}

	offset := (param.Page - 1) * param.Perpage
	dbQuery := r.db.Model(&AppModul{}).
		Preload("Features").
		Order("created_at DESC")
	if param.Search != "" {
		search := "%" + param.Search + "%"
		dbQuery = dbQuery.Where(
			"name ILIKE ? OR code ILIKE ?",
			search, search,
		)
	}
	dbQuery.Count(&total)
	dbQuery.Limit(param.Perpage).Offset(offset).Find(&data)

	totalPage := int((total + int64(param.Perpage) - 1) / int64(param.Perpage))

	return &AppModulIndex{
		Data:      data,
		Total:     total,
		TotalPage: totalPage,
		Page:      param.Page,
		PerPage:   param.Perpage,
	}, nil
}
func (r *repository) ModulCreate(req ModulRequest) (*AppModul, error) {
	modulId := uuid.NewString()
	timeNow := time.Now()
	err := r.db.Transaction(func(tx *gorm.DB) error {
		create := AppModul{
			Id:          modulId,
			AppId:       req.AppId,
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
			Status:      req.Status,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}
		if err := tx.Create(&create).Error; err != nil {
			return err
		}
		for _, item := range req.Features {
			Id := uuid.NewString()
			record := AppModulFeature{
				Id:          Id,
				AppId:       req.AppId,
				AppModulId:  modulId,
				Code:        item.Code,
				Name:        item.Name,
				Description: item.Description,
				Permission:  item.Permission,
				Status:      item.Status,
				CreatedAt:   timeNow,
				UpdatedAt:   timeNow,
			}
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("failed to create modul")
	}
	res, _ := r.ModulRead(modulId)
	return res, nil
}
func (r *repository) ModulUpdate(req ModulRequest, Id string) (*AppModul, error) {
	updateData := map[string]interface{}{
		"code":        req.Code,
		"name":        req.Name,
		"description": req.Description,
		"status":      req.Status,
		"updated_at":  time.Now(),
	}
	result := r.db.Model(&AppModul{}).Where("id = ?", Id).Updates(updateData)
	if result.Error != nil {
		return nil, errors.New("update failed")
	}
	res, _ := r.ModulRead(Id)
	return res, nil
}
func (r *repository) ModulUpdateStatus(req ModulRequest, Id string) (*AppModul, error) {
	updateData := map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	}
	result := r.db.Model(&AppModul{}).Where("id = ?", Id).Updates(updateData)
	if result.Error != nil {
		return nil, errors.New("update failed")
	}
	res, _ := r.ModulRead(Id)
	return res, nil
}
func (r *repository) ModulDelete(id string) error {
	delete := r.db.Where("id = ?", id).Delete(&AppModul{})
	if delete.Error != nil {
		return errors.New("data cannot be deleted")
	}

	return nil
}
func (r *repository) ModulRead(id string) (*AppModul, error) {
	var data AppModul
	check := r.db.Model(AppModul{}).Preload("Features").Where("id = ?", id).Find(&data)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("id not found")
	}
	return &data, nil
}

func (r *repository) FeatureCreate(req FeatureRequest) (*AppModulFeature, error) {
	id := uuid.NewString()
	timeNow := time.Now()
	plugin := AppModulFeature{
		Id:          id,
		AppId:       req.AppId,
		AppModulId:  req.AppModulId,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Permission:  req.Permission,
		Status:      req.Status,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	result := r.db.Create(&plugin)

	if result.Error != nil {
		return nil, result.Error
	}
	res, _ := r.FeatureRead(id)
	return res, nil
}
func (r *repository) FeatureUpdate(req FeatureRequest, Id string) (*AppModulFeature, error) {
	updateData := map[string]interface{}{
		"code":        req.Code,
		"name":        req.Name,
		"description": req.Description,
		"permission":  req.Permission,
		"status":      req.Status,
		"updated_at":  time.Now(),
	}
	result := r.db.Model(&AppModulFeature{}).Where("id = ?", Id).Updates(updateData)
	if result.Error != nil {
		return nil, errors.New("update failed")
	}
	res, _ := r.FeatureRead(Id)
	return res, nil
}
func (r *repository) FeatureUpdateStatus(req FeatureRequest, Id string) (*AppModulFeature, error) {
	updateData := map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	}
	result := r.db.Model(&AppModulFeature{}).Where("id = ?", Id).Updates(updateData)
	if result.Error != nil {
		return nil, errors.New("update failed")
	}
	res, _ := r.FeatureRead(Id)
	return res, nil
}
func (r *repository) FeatureDelete(id string) error {
	delete := r.db.Where("id = ?", id).Delete(&AppModulFeature{})
	if delete.Error != nil {
		return errors.New("data cannot be deleted")
	}

	return nil
}
func (r *repository) FeatureRead(id string) (*AppModulFeature, error) {
	var data AppModulFeature
	check := r.db.Model(AppModulFeature{}).Where("id = ?", id).Find(&data)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("id not found")
	}
	return &data, nil
}
func (r *repository) FeatureBulkDelete(ids []string) (int, int, error) {

	delete := r.db.Where("id IN ?", ids).Delete(&AppModulFeature{})
	if delete.Error != nil {
		return 0, 0, errors.New("data cannot be deleted")
	}

	success := int(delete.RowsAffected)
	failed := len(ids) - success

	return success, failed, nil
}

func (r *repository) CheckAppId(appId string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM apps WHERE id = ?)`
	if err := r.db.Raw(query, appId).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}
func (r *repository) CheckAppCode(code string, id ...string) bool {
	var exists bool
	if id != nil {
		query := `SELECT EXISTS (SELECT 1 FROM apps WHERE code = ? AND id != ?)`
		if err := r.db.Raw(query, code, id).Scan(&exists).Error; err != nil {
			return false
		}
	} else {
		query := `SELECT EXISTS (SELECT 1 FROM apps WHERE code = ?)`
		if err := r.db.Raw(query, code).Scan(&exists).Error; err != nil {
			return false
		}
	}
	return exists
}
func (r *repository) CheckPluginCode(code string, id ...string) bool {
	var exists bool
	if id != nil {
		query := `SELECT EXISTS (SELECT 1 FROM app_plugins WHERE code = ? AND id != ?)`
		if err := r.db.Raw(query, code, id).Scan(&exists).Error; err != nil {
			return false
		}
	} else {
		query := `SELECT EXISTS (SELECT 1 FROM app_plugins WHERE code = ?)`
		if err := r.db.Raw(query, code).Scan(&exists).Error; err != nil {
			return false
		}
	}
	return exists
}

func (r *repository) CheckAppModulId(appModulId string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM app_moduls WHERE id = ?)`
	if err := r.db.Raw(query, appModulId).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}
