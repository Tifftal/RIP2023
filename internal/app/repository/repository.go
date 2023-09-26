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

	err := repository.db.First(sample, "Id_sample = ?", id).Error // find alpinist with code D42
	if err != nil {
		return nil, err
	}

	return sample, nil
}

func (repository *Repository) GetAllSamples() ([]ds.Samples, error) {
	sample := []ds.Samples{}
	err := repository.db.Find(&sample).Error
	if err != nil {
		return nil, err
	}

	return sample, nil
}
