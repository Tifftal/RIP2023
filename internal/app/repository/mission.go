package repository

import (
	"MSRM/internal/app/ds"
	"errors"
	"time"

	"gorm.io/gorm"
)

func (repository *Repository) GetAllMissions(user_id int) ([]ds.Missions, error) {
	missions := []ds.Missions{}
	var user ds.Users
	err := repository.db.Table("users").Where("Id_user = ? AND (Role = 'Moderator' OR Role = 'User')", user_id).First(&user).Error
	if err != nil {
		return nil, errors.New("Недостаточно прав для просмотра миссии")
	}

	if repository.db.Table("users").Where("Id_user = ? AND Role = 'User'", user_id).First(&user).RowsAffected > 0 {
		err := repository.db.Where("user_id = ?", user_id).Order("formation_date ASC").Find(&missions).Error
		if err != nil {
			return nil, err
		}
	}

	if repository.db.Table("users").Where("Id_user = ? AND Role = 'Moderator'", user_id).First(&user).RowsAffected > 0 {
		err := repository.db.Where("moderator_id = ?", user_id).Order("formation_date ASC").Find(&missions).Error
		if err != nil {
			return nil, err
		}
	}

	return missions, nil
}

func (repository *Repository) GetAllMissionsByDateRange(startDate, endDate time.Time, user_id int) ([]ds.Missions, error) {
	missions := []ds.Missions{}
	var user ds.Users
	err := repository.db.Table("users").Where("Id_user = ? AND (Role = 'Moderator' OR Role = 'User')", user_id).First(&user).Error
	if err != nil {
		return nil, errors.New("Недостаточно прав для просмотра миссии")
	}

	// Для пользователя с ролью 'User'
	if err := repository.db.Table("users").Where("Id_user = ? AND Role = 'User'", user_id).First(&user).Error; err == nil {
		err := repository.db.Where("formation_date BETWEEN ? AND ? AND user_id = ?", startDate, endDate, user_id).Order("formation_date ASC").Find(&missions).Error
		if err != nil {
			return nil, err
		}
	}

	// Для пользователя с ролью 'Moderator'
	if err := repository.db.Table("users").Where("Id_user = ? AND Role = 'Moderator'", user_id).First(&user).Error; err == nil {
		err := repository.db.Where("formation_date BETWEEN ? AND ? AND moderator_id = ?", startDate, endDate, user_id).Order("formation_date ASC").Find(&missions).Error
		if err != nil {
			return nil, err
		}
	}

	return missions, nil
}

func (repository *Repository) GetMissionByID(id int) (*ds.Missions, error) {
	mission := &ds.Missions{}

	err := repository.db.First(mission, "Id_mission = ?", id).Error
	if err != nil {
		return nil, err
	}

	return mission, nil
}

func (r *Repository) UpdateMission(mission *ds.Missions, id int, user_id int) error {
	// Проверяем, что пользователь авторизирован
	var user ds.Users
	err := r.db.Table("users").Where("Id_user = ? AND Role = 'Moderator'", user_id).First(&user).Error
	if err != nil {
		return errors.New("Недостаточно прав для редактирования миссии")
	}

	// Проверяем, что миссия существует и принадлежит пользователю
	err = r.db.Table("missions").
		Where("id_mission = ? AND moderator_id = ?", id, user_id).
		First(&ds.Missions{}).Error

	if err != nil {
		return errors.New("Эта миссия не принадлежит вам")
	}

	// Теперь обновляем миссию
	updateErr := r.db.Where("Id_mission = ? AND moderator_id = ?", id, user_id).Updates(&mission).Error
	return updateErr

}

func (r *Repository) DeleteMissionByID(id int, user_id int) error {
	var user ds.Users
	err := r.db.Table("users").Where("Id_user = ? AND Role = 'User'", user_id).First(&user).Error
	if err != nil {
		return errors.New("Недостаточно прав для удаления миссии")
	}
	if err := r.db.Exec("UPDATE missions SET mission_status='Deleted' WHERE Id_mission= ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (repository *Repository) GetMissioninDetailByID(id int, user_id int) (*ds.Missions, []ds.Samples, error) {
	mission := &ds.Missions{}
	samples := []ds.Samples{}
	var user ds.Users

	err := repository.db.Table("users").
		Where("Id_user = ? AND (Role = 'Moderator' OR Role = 'User')", user_id).
		First(&user).Error

	if err != nil {
		return nil, nil, errors.New("Недостаточно прав для редактирования миссии")
	}

	// Проверка, принадлежит ли миссия пользователю или он модератор
	err = repository.db.Table("missions").
		Where("id_mission = ? AND (user_id = ? OR moderator_id = ?)", id, user_id, user_id).
		First(&mission).Error

	if err != nil {
		return nil, nil, errors.New("Эта миссия не принадлежит вам")
	}

	// Выборка связанных образцов
	err = repository.db.
		Joins("JOIN mission_samples ON missions.id_mission = mission_samples.mission_id").
		Joins("JOIN samples ON mission_samples.sample_id = samples.id_sample").
		Where("missions.id_mission = ?", id).
		Table("missions").
		Select("missions.*, samples.*").
		Find(&samples).Error

	if err != nil {
		return nil, nil, err
	}

	return mission, samples, nil
}

func (r *Repository) UpdateMissionStatusByUser(id int, newStatus string, user_id int) error {
	var user ds.Users
	err := r.db.Table("users").Where("Id_user = ? AND Role = 'User'", user_id).First(&user).Error
	if err != nil {
		return errors.New("Недостаточно прав для редактирования миссии")
	}

	var currentStatus string
	err = r.db.Table("missions").Where("Id_mission = ? AND user_id = ?", id, user_id).Pluck("mission_status", &currentStatus).Error
	if err != nil {
		return errors.New("Не удалось получить текущий статус миссии")
	}

	// Проверяем, что новый статус допустим
	allowedStatus := map[string]bool{
		"Draft":                 true,
		"Deleted":               true,
		"Awaiting confirmation": true,
	}

	if !allowedStatus[newStatus] {
		return errors.New("Неправильный статус миссии")
	}

	switch currentStatus {
	case "Draft":
		if !(newStatus == "Awaiting confirmation" || newStatus == "Deleted") {
			return errors.New("Нельзя изменить статус Draft на " + newStatus)
		}
	case "Deleted":
		return errors.New("Нельзя изменить статус Deleted")
	}

	// Обновляем статус миссии
	updateErr := r.db.Model(&ds.Missions{}).
		Where("Id_mission = ? AND user_id = ?", id, user_id).
		Update("mission_status", newStatus).
		Error

	return updateErr
}

func (r *Repository) UpdateMissionStatusByModerator(id int, newStatus string, user_id int) error {
	var user ds.Users
	err := r.db.Table("users").Where("Id_user = ? AND Role = 'Moderator'", user_id).First(&user).Error
	if err != nil {
		return errors.New("Недостаточно прав для редактирования миссии")
	}

	var currentStatus string
	err = r.db.Table("missions").Where("Id_mission = ? AND moderator_id = ?", id, user_id).Pluck("mission_status", &currentStatus).Error
	if err != nil {
		return errors.New("Не удалось получить текущий статус миссии")
	}

	// Проверяем, что новый статус допустим
	allowedStatus := map[string]bool{
		"Completed": true,
		"Rejected":  true,
		"At work":   true,
	}

	if !allowedStatus[newStatus] {
		return errors.New("Неправильный статус миссии")
	}

	switch currentStatus {
	case "Awaiting confirmation":
		if !(newStatus == "At work" || newStatus == "Rejected") {
			return errors.New("Нельзя изменить статус Awaiting confirmation на " + newStatus)
		}
	case "Draft":
		return errors.New("Нельзя изменить статус Draft")
	case "At work":
		if !(newStatus == "Completed" || newStatus == "Rejected") {
			return errors.New("Нельзя изменить статус At work на " + newStatus)
		}
	}

	if !allowedStatus[newStatus] {
		return errors.New("Неправильный статус миссии")
	}

	// Обновляем статус миссии
	updateErr := r.db.Model(&ds.Missions{}).
		Where("Id_mission = ? AND moderator_id = ?", id, user_id).
		Update("mission_status", newStatus).
		Error

	return updateErr
}

func (repository *Repository) RemoveSampleFromLastDraftMission(sampleID int, user_id int) (*ds.Missions, []ds.Samples, error) {
	var user ds.Users
	err := repository.db.Table("users").Where("Id_user = ? AND Role = 'User'", user_id).First(&user).Error
	if err != nil {
		return nil, nil, errors.New("Недостаточно прав для редактирования миссии")
	}
	// Находим последнюю миссию с mission_status = "Draft"
	var lastDraftMission ds.Missions
	dbErr := repository.db.
		Order("formation_date desc").
		Where("mission_status = ? AND user_id = ?", "Draft", user_id).
		First(&lastDraftMission).
		Error

	if dbErr != nil && !errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return nil, nil, dbErr
	}

	// Если миссии с mission_status = "Draft" нет, возвращаем ошибку
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return nil, nil, errors.New("Миссия со статусом Draft не найдена")
	}

	// Удаляем образец из миссии
	if err := repository.db.Exec("DELETE FROM mission_samples WHERE mission_id = ? AND sample_id = ?", lastDraftMission.Id_mission, sampleID).Error; err != nil {
		return nil, nil, err
	}

	// Получаем все образцы в миссии после удаления
	var samples []ds.Samples
	removeErr := repository.db.
		Joins("JOIN mission_samples ON missions.id_mission = mission_samples.mission_id").
		Joins("JOIN samples ON mission_samples.sample_id = samples.id_sample").
		Where("missions.id_mission = ?", lastDraftMission.Id_mission).
		Table("missions").
		Select("missions.*, samples.*").
		Find(&samples).
		Error

	if removeErr != nil {
		return nil, nil, removeErr
	}

	return &lastDraftMission, samples, nil
}
