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

func (repository *Repository) AddSample(sample *ds.Samples) (bool, error) {
	err := repository.db.Create(sample).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) DeleteSampleByID(id int) error {
	if err := r.db.Exec("UPDATE samples SET sample_status='Deleted' WHERE Id_sample= ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) ReturnSampleByID(id int) error {
	if err := r.db.Exec("UPDATE samples SET sample_status='Active' WHERE Id_sample= ?", id).Error; err != nil {
		return err
	}
	return nil
}
