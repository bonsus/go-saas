package media

import (
	"errors"
	"time"

	"github.com/bonsus/go-saas/internal/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaRepository interface {
	// Create(Option Option) (*Option, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}
func (r *repository) Create(req Medias) (*Medias, error) {
	Id := uuid.NewString()

	tx := r.db.Begin()

	query := `
		INSERT INTO medias (id, name, description, alt, type, status, created_at, updated_at)
		VALUES (?,?,?,?,?,?,?,?)
	`
	timeNow := time.Now()
	result := tx.Exec(query, Id, req.Name, req.Description, req.Alt, req.Type, "public", timeNow, timeNow)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	for _, file := range req.Files {
		fileId := uuid.NewString()
		query := `INSERT INTO media_files (id, media_id, file, width, height, filesize) VALUES (?,?,?,?,?,?)`
		result := tx.Exec(query, fileId, Id, file.File, file.Width, file.Height, file.Filesize)
		if result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	media, err := r.Read(Id)
	if err != nil {
		return nil, errors.New("data not found")
	}
	return media, nil
}

func (r *repository) check(id string) bool {
	var exists bool
	r.db.Raw("SELECT EXISTS (SELECT 1 FROM medias WHERE id = ?)", id).Scan(&exists)
	return exists
}
func (r *repository) Update(req updateRequest, Id string) (*Medias, error) {

	query := `
		UPDATE medias SET name = ?, description = ?, alt = ?, updated_at = ?
		WHERE id = ?
	`
	result := r.db.Exec(query, req.Name, req.Description, req.Alt, time.Now(), Id)
	if result.Error != nil {
		return nil, result.Error
	}
	media, err := r.Read(Id)
	if err != nil {
		return nil, errors.New("data not found")
	}
	return media, nil
}

func (r *repository) Read(Id string) (*Medias, error) {
	var media Medias
	var cfg = config.GetConfig()
	err := r.db.Table("medias").Preload("Files").Where("id =?", Id).First(&media).Error
	if err != nil || media.Id == "" {
		return nil, errors.New("data not found")
	}
	for i, f := range media.Files {
		media.Files[i].Url = cfg.S3.Domain + f.File
	}
	return &media, nil
}

func (r *repository) Index(param ParamIndex) (result *MediaIndex, err error) {
	var data []Medias
	var total int64

	if param.Page < 1 {
		param.Page = 1
	}
	if param.Perpage < 1 {
		param.Perpage = 20
	}
	offset := (param.Page - 1) * param.Perpage
	dbQuery := r.db.Model(&Medias{}).Preload("Files").Where("medias.deleted_at IS NULL").Group("medias.id")
	if len(param.Search) > 0 {
		search := "%" + param.Search + "%"
		dbQuery = dbQuery.Joins("JOIN media_files ON media_files.media_id = medias.id")
		dbQuery = dbQuery.Where("medias.name ILIKE ? OR medias.description ILIKE ? OR medias.alt ILIKE ? or media_files.file ILIKE ?", search, search, search, search)
	}
	if len(param.Status) > 0 {
		dbQuery = dbQuery.Where("media.status = ?", param.Status)
	}
	dbQuery.Count(&total)
	dbQuery.Limit(param.Perpage).Offset(offset).Find(&data)

	totalPage := int((total + int64(param.Perpage) - 1) / int64(param.Perpage))

	cfg := config.GetConfig()
	for _, media := range data {
		for i, f := range media.Files {
			media.Files[i].Url = cfg.S3.Domain + f.File
		}
	}
	return &MediaIndex{
		Data:      data,
		Total:     total,
		TotalPage: totalPage,
		Page:      param.Page,
		PerPage:   param.Perpage,
	}, nil
}

func (r *repository) Delete(id string) (err error) {

	tx := r.db.Begin()
	if err := tx.Exec("DELETE FROM media_files WHERE media_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}
	result := tx.Exec("DELETE FROM medias WHERE id = ?", id)
	if result.Error != nil {
		tx.Rollback()
		return errors.New("cannot delete customer: it is still assigned to other data")
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("no customer deleted")
	}
	return tx.Commit().Error
}
