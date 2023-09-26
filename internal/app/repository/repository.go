package repository

import (
	"MSRM/internal/app/ds"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (repository *Repository) GetSampleByID(id int) (*ds.Samples, error) {
	sample := &ds.Samples{}

	err := repository.db.First(sample, "Id_sample = ?", id).Error
	if err != nil {
		return nil, err
	}

	return sample, nil
}

func (repository *Repository) GetAllSamples() ([]ds.Samples, error) {
	sample := []ds.Samples{}
	err := repository.db.Order("Sample_status ASC").Order("Id_sample ASC").Find(&sample).Error
	if err != nil {
		return nil, err
	}

	return sample, nil
}

func (repository *Repository) GetSampleByName(name string) ([]ds.Samples, error) {
	var samples []ds.Samples
	err := repository.db.Where("Name LIKE ?", "%"+name+"%").Order("Sample_status ASC").Order("Id_sample ASC").Find(&samples).Error
	if err != nil {
		return nil, err
	}
	return samples, nil
}
