package repository

import (
	"MSRM/internal/app/ds"
	"errors"

	"gorm.io/gorm"
)

func (repository *Repository) GetAllMissions() ([]ds.Missions, error) {
	mission := []ds.Missions{}
	err := repository.db.Find(&mission).Error
	if err != nil {
		return nil, err
	}

	return mission, nil
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
	err := r.db.Table("users").Where("Id_user = ?", user_id).First(&user).Error
	if err != nil {
		return errors.New("Чтобы редактировать миссию, нужно авторизироваться")
	}
	updateErr := r.db.Where("Id_mission = ?", id).Updates(&mission).Error
	return updateErr
}

func (r *Repository) DeleteMissionByID(id int, user_id int) error {
	var user ds.Users
	err := r.db.Table("users").Where("Id_user = ? AND Role = 'Moderator'", user_id).First(&user).Error
	if err != nil {
		return errors.New("Недостаточно прав для удаления миссии")
	}
	if err := r.db.Exec("UPDATE missions SET mission_status='Deleted' WHERE Id_mission= ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (repository *Repository) GetMissioninDetailByID(id int) (*ds.Missions, []ds.Samples, error) {
	mission := &ds.Missions{}
	samples := []ds.Samples{}

	// Retrieve mission details
	err := repository.db.First(mission, "Id_mission = ?", id).Error
	if err != nil {
		return nil, nil, err
	}

	// Retrieve associated samples
	err = repository.db.
		Joins("JOIN mission_samples ON missions.id_mission = mission_samples.id_mission").
		Joins("JOIN samples ON mission_samples.id_sample = samples.id_sample").
		Where("missions.id_mission = ?", id).
		Table("missions"). // Add this line to specify the table name
		Select("missions.*, samples.*").
		Find(&samples).Error

	if err != nil {
		return nil, nil, err
	}

	return mission, samples, nil
}

// func (repository *Repository) GetMissionByUserID(id int) ([]ds.Missions, error) {
// 	mission := []ds.Missions{}
// 	err := repository.db.Where("User_id=?", id).Find(&mission).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return mission, nil
// }

// func (repository *Repository) GetMissionByModeratorID(id int) ([]ds.Missions, error) {
// 	mission := []ds.Missions{}
// 	err := repository.db.Where("Moderator_id=?", id).Find(&mission).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return mission, nil
// }

// func (repository *Repository) GetMissionByStatus(status string) ([]ds.Missions, error) {
// 	mission := []ds.Missions{}
// 	err := repository.db.Where("Mission_status=?", status).Find(&mission).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return mission, nil
// }

func (r *Repository) UpdateMissionStatusByUser(id int, newStatus string, user_id int) error {
	var user ds.Users
	err := r.db.Table("users").Where("Id_user = ? AND Role = 'User'", user_id).First(&user).Error
	if err != nil {
		return errors.New("Недостаточно прав для редактирования миссии")
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

	// Обновляем статус миссии
	updateErr := r.db.Model(&ds.Missions{}).
		Where("Id_mission = ?", id).
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
	// Проверяем, что новый статус допустим
	allowedStatus := map[string]bool{
		"Completed": true,
		"Rejected":  true,
		"At work":   true,
	}

	if !allowedStatus[newStatus] {
		return errors.New("Неправильный статус миссии")
	}

	// Обновляем статус миссии
	updateErr := r.db.Model(&ds.Missions{}).
		Where("Id_mission = ?", id).
		Update("mission_status", newStatus).
		Error

	return updateErr
}

func (repository *Repository) RemoveSampleFromMission(missionID, sampleID uint, user_id int) (*ds.Missions, []ds.Samples, error) {
	// Проверяем, существует ли миссия с указанным ID
	var mission ds.Missions
	if err := repository.db.Where("mission_status != ?", "Draft").First(&mission, missionID).Error; err != nil {
		return nil, nil, errors.New("Нельзя удалять миссию со статусом Draft")
	}

	// Проверяем, что user_id совпадает с user_id из миссии
	if mission.Moderator_id != user_id {
		return nil, nil, errors.New("Недостаточно прав для редактирования этой миссии")
	}

	var user ds.Users
	err := repository.db.Table("users").Where("Id_user = ? AND Role = 'Moderator'", user_id).First(&user).Error
	if err != nil {
		return nil, nil, errors.New("Недостаточно прав для редактирования миссии")
	}

	// Удаляем запись об образце из таблицы mission_samples
	if err := repository.db.
		Where("mission_id = ? AND sample_id = ?", missionID, sampleID).
		Delete(&ds.Mission_samples{}).Error; err != nil {
		return nil, nil, err
	}

	// Получаем все образцы в миссии после удаления
	var samples []ds.Samples
	removeErr := repository.db.
		Joins("JOIN mission_samples ON missions.id_mission = mission_samples.mission_id").
		Joins("JOIN samples ON mission_samples.sample_id = samples.id_sample").
		Where("missions.id_mission = ?", missionID).
		Table("missions").
		Select("missions.*, samples.*").
		Find(&samples).
		Error

	if err != nil {
		return nil, nil, removeErr
	}

	return &mission, samples, nil
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
		Where("mission_status = ?", "Draft").
		First(&lastDraftMission).
		Error

	if dbErr != nil && !errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return nil, nil, dbErr
	}

	// Если миссии с mission_status = "Draft" нет, возвращаем ошибку
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return nil, nil, errors.New("Миссия со статусом Draft не найдена")
	}

	// Получаем миссию по lastDraftMission.Id_mission
	var missionWithUserID ds.Missions
	if err := repository.db.First(&missionWithUserID, lastDraftMission.Id_mission).Error; err != nil {
		return nil, nil, err
	}

	// Сравниваем user_id из миссии с переданным user_id
	if missionWithUserID.User_id != user_id {
		return nil, nil, errors.New("Недостаточно прав для удаления образца из миссии")
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
