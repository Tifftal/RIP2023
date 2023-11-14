package repository

import (
	"MSRM/internal/app/ds"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

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

// func (repository *Repository) GetSampleByName(name string) ([]ds.Samples, error) {
// 	var samples []ds.Samples
// 	err := repository.db.Where("Name LIKE ?", "%"+name+"%").Order("Sample_status ASC").Order("Id_sample ASC").Find(&samples).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return samples, nil
// }

func (r *Repository) DeleteSampleByID(id, user_id int) error {
	// Получаем информацию о пользователе по user_id
	var user ds.Users
	err := r.db.Table("users").Where("Id_user = ? AND Role = 'Moderator'", user_id).First(&user).Error
	if err != nil {
		return errors.New("недостаточно прав для удаления образца")
	}

	// Выполняем обновление только если проверка пройдена
	if err := r.db.Exec("UPDATE samples SET sample_status='Deleted' WHERE Id_sample = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateSample(sample *ds.Samples, user_id int) error {
	var user ds.Users
	err := r.db.Table("users").Where("Id_user = ? AND Role = 'Moderator'", user_id).First(&user).Error
	if err != nil {
		return errors.New("Недостаточно прав для создания образца")
	}

	// Use a different variable name for the Create error
	createErr := r.db.Create(&sample).Error

	return createErr
}

func (r *Repository) UpdateSample(sample *ds.Samples, user_id int) error {
	var user ds.Users
	err := r.db.Table("users").Where("Id_user = ? AND Role = 'Moderator'", user_id).First(&user).Error
	if err != nil {
		return errors.New("Недостаточно прав для редактирования образца")
	}

	updateErr := r.db.Where("Id_sample = ?", sample.Id_sample).Updates(&sample).Error

	return updateErr
}

func (r *Repository) AddSampleToLastDraftMission(sampleID int, user_id int) (*ds.Missions, []ds.Samples, error) {
	var user ds.Users
	err := r.db.Table("users").Where("Id_user = ? AND Role = 'User'", user_id).First(&user).Error
	if err != nil {
		return nil, nil, errors.New("Необходимо авторизироваться для добавления образца в миссию")
	}
	// Находим последнюю миссию с mission_status = "Draft"
	var lastDraftMission ds.Missions
	dbErr := r.db.
		Table("missions").
		Order("formation_date desc").
		Where("user_id = ? AND mission_status = 'Draft'", user_id).
		First(&lastDraftMission).
		Error

	if dbErr != nil && !errors.Is(dbErr, gorm.ErrRecordNotFound) {
		fmt.Println("не вот тут")
		return nil, nil, dbErr
	}
	fmt.Println(lastDraftMission)

	// Если миссии с mission_status = "Draft" нет, создаем новую
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		currentTime := time.Now()
		lastDraftMission = ds.Missions{
			User_id:        user_id,
			Moderator_id:   2,
			Mission_status: "Draft",
			Name:           "NewDraftMission",
			Creation_date:  currentTime,
			Formation_date: currentTime,
		}
		if err := r.db.Create(&lastDraftMission).Error; err != nil {
			return nil, nil, err
		}
	}
	// Получаем образец из базы данных по его идентификатору
	var newSample ds.Samples
	if err := r.db.First(&newSample, sampleID).Error; err != nil {
		return nil, nil, err
	}

	// Добавляем образец в миссию
	if err := r.db.Create(&ds.Mission_samples{
		Mission_id: lastDraftMission.Id_mission,
		Sample_id:  newSample.Id_sample,
	}).Error; err != nil {
		// Проверяем, является ли ошибка уникальным ключом
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			// Здесь обрабатываем случай дубликата ключа, если это произошло
			return nil, nil, errors.New("Образец уже добавлен в миссию")
		}
		return nil, nil, err
	}

	// Получаем все образцы в миссии
	var samples []ds.Samples
	addErr := r.db.
		Joins("JOIN mission_samples ON missions.id_mission = mission_samples.mission_id").
		Joins("JOIN samples ON mission_samples.sample_id = samples.id_sample").
		Where("missions.id_mission = ?", lastDraftMission.Id_mission).
		Table("missions").
		Select("missions.*, samples.*").
		Find(&samples).
		Error

	if addErr != nil {
		fmt.Println("вот тут ошибка")
		return nil, nil, addErr
	}

	return &lastDraftMission, samples, nil

}

// func (repository *Repository) GetAllSamplesOrderByType() ([]ds.Samples, error) {
// 	sample := []ds.Samples{}
// 	err := repository.db.Order("Type ASC").Order("Id_sample ASC").Find(&sample).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return sample, nil
// }

// func (repository *Repository) GetAllSamplesOrderByDate() ([]ds.Samples, error) {
// 	sample := []ds.Samples{}
// 	err := repository.db.Order("Date_Sealed ASC").Order("Id_sample ASC").Find(&sample).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return sample, nil
// }

// func (repository *Repository) GetAllSamplesStatusActive() ([]ds.Samples, error) {
// 	sample := []ds.Samples{}
// 	err := repository.db.Where("sample_status='Active'").Order("Id_sample ASC").Find(&sample).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return sample, nil
// }

// func (repository *Repository) GetAllSamplesStatusDeleted() ([]ds.Samples, error) {
// 	sample := []ds.Samples{}
// 	err := repository.db.Where("sample_status='Deleted'").Order("Id_sample ASC").Find(&sample).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return sample, nil
// }
