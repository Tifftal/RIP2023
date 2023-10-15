package repository

import "MSRM/internal/app/ds"

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

func (r *Repository) DeleteSampleByID(id int) error {
	if err := r.db.Exec("UPDATE samples SET sample_status='Deleted' WHERE Id_sample= ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateSample(sample *ds.Samples) error {
	err := r.db.Create(&sample).Error

	return err
}

func (r *Repository) UpdateSample(sample *ds.Samples) error {
	err := r.db.Where("Id_sample = ?", sample.Id_sample).Updates(&sample).Error

	return err
}

func (repository *Repository) GetAllSamplesOrderByType() ([]ds.Samples, error) {
	sample := []ds.Samples{}
	err := repository.db.Order("Type ASC").Order("Id_sample ASC").Find(&sample).Error
	if err != nil {
		return nil, err
	}

	return sample, nil
}

func (repository *Repository) GetAllSamplesOrderByDate() ([]ds.Samples, error) {
	sample := []ds.Samples{}
	err := repository.db.Order("Date_Sealed ASC").Order("Id_sample ASC").Find(&sample).Error
	if err != nil {
		return nil, err
	}

	return sample, nil
}

func (repository *Repository) GetAllSamplesStatusActive() ([]ds.Samples, error) {
	sample := []ds.Samples{}
	err := repository.db.Where("sample_status='Active'").Order("Id_sample ASC").Find(&sample).Error
	if err != nil {
		return nil, err
	}

	return sample, nil
}

func (repository *Repository) GetAllSamplesStatusDeleted() ([]ds.Samples, error) {
	sample := []ds.Samples{}
	err := repository.db.Where("sample_status='Deleted'").Order("Id_sample ASC").Find(&sample).Error
	if err != nil {
		return nil, err
	}

	return sample, nil
}
