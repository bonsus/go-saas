package country

import (
	"gorm.io/gorm"
)

type Repository interface {
	Provinces(param Params) ([]Province, error)
	Cities(param Params) ([]City, error)
	Districts(param Params) ([]District, error)
	Zips(param Params) ([]Zip, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

func (r *repository) Countries(param Params) ([]Country, error) {
	var countries []Country

	query := `
		SELECT country
		FROM countries 
		WHERE country <> ''
		GROUP BY country
		ORDER BY country ASC
	`

	err := r.db.Raw(query).Scan(&countries).Error
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *repository) Provinces(param Params) ([]Province, error) {
	var provinces []Province

	query := `
		SELECT country, province
		FROM countries 
		WHERE province <> ''
		GROUP BY country, province
		ORDER BY country, province ASC
	`

	err := r.db.Raw(query).Scan(&provinces).Error
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

func (r *repository) Cities(param Params) (*Cities, error) {
	var cities []City

	query := `
		SELECT country, province, city
		FROM countries 
		WHERE province = ? AND city IS NOT NULL AND city <> ''
		GROUP BY country, province, city
		ORDER BY country, province, city ASC
	`

	err := r.db.Raw(query, param.Province).Scan(&cities).Error
	if err != nil {
		return nil, err
	}

	return &Cities{
		Province: param.Province,
		Cities:   cities,
	}, nil
}
func (r *repository) Districts(param Params) (*Districts, error) {
	var districts []District

	query := `
		SELECT country, province, city, district
		FROM countries 
		WHERE province = ? AND city = ? AND district IS NOT NULL AND district <> ''
		GROUP BY country, province, city, district
		ORDER BY country, province, city, district ASC
	`

	err := r.db.Raw(query, param.Province, param.City).Scan(&districts).Error
	if err != nil {
		return nil, err
	}

	return &Districts{
		Province:  param.Province,
		City:      param.City,
		Districts: districts,
	}, nil
}

func (r *repository) Zips(param Params) (*Zips, error) {
	var zips []Zip

	query := `
		SELECT country, province, city, district,zip
		FROM countries 
		WHERE province = ? AND city = ? AND district = ? AND zip IS NOT NULL AND zip <> ''
		GROUP BY country, province, city, district, zip
		ORDER BY country, province, city, district, zip ASC
	`

	err := r.db.Raw(query, param.Province, param.City, param.District).Scan(&zips).Error
	if err != nil {
		return nil, err
	}

	return &Zips{
		Province: param.Province,
		City:     param.City,
		District: param.District,
		Zips:     zips,
	}, nil
}

func (r *repository) Search(param Params) ([]Search, error) {
	var searches []Search
	search := "%" + param.Search + "%"
	query := `
		(
			SELECT country, province, city, district
			FROM countries
			WHERE COALESCE(district, '') ILIKE ?
				AND district <> ''
			GROUP BY country, province, city, district
			ORDER BY district ASC
			LIMIT 20
		)
		UNION ALL
		(
			SELECT country, province, city, district
			FROM countries
			WHERE COALESCE(city, '') ILIKE ?
				AND district <> ''
				AND NOT EXISTS (
					SELECT 1 FROM countries WHERE COALESCE(district, '') ILIKE ?
				)
			GROUP BY country, province, city, district
			ORDER BY city ASC
			LIMIT 20
		)
	`

	err := r.db.Raw(query, search, search, search).Scan(&searches).Error
	if err != nil {
		return nil, err
	}

	return searches, nil
}
