package repository

import (
	"MSRM/internal/app/ds"
	"errors"
	"fmt"
	"time"

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

func (r *Repository) DeleteMissionByID(id int) error {
	if err := r.db.Exec("UPDATE missions SET mission_status='Deleted' WHERE Id_mission= ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateMission(mission *ds.Missions) error {
	err := r.db.Where("Id_mission = ?", mission.Id_mission).Updates(&mission).Error

	return err
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

func (repository *Repository) GetMissionByUserID(id int) ([]ds.Missions, error) {
	mission := []ds.Missions{}
	err := repository.db.Where("User_id=?", id).Find(&mission).Error
	if err != nil {
		return nil, err
	}

	return mission, nil
}

func (repository *Repository) GetMissionByModeratorID(id int) ([]ds.Missions, error) {
	mission := []ds.Missions{}
	err := repository.db.Where("Moderator_id=?", id).Find(&mission).Error
	if err != nil {
		return nil, err
	}

	return mission, nil
}

func (repository *Repository) GetMissionByStatus(status string) ([]ds.Missions, error) {
	mission := []ds.Missions{}
	err := repository.db.Where("Mission_status=?", status).Find(&mission).Error
	if err != nil {
		return nil, err
	}

	return mission, nil
}

func (repository *Repository) AddSampleToLastDraftMission(sampleID int) (*ds.Missions, []ds.Samples, error) {
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
	fmt.Println(lastDraftMission)

	// Если миссии с mission_status = "Draft" нет, создаем новую
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		currentTime := time.Now()
		lastDraftMission = ds.Missions{
			Mission_status: "Draft",
			Name:           "NewDraftMission",
			Creation_date:  currentTime,
			Formation_date: currentTime,
		}
		if err := repository.db.Create(&lastDraftMission).Error; err != nil {
			return nil, nil, err
		}
	}
	// Получаем образец из базы данных по его идентификатору
	var newSample ds.Samples
	if err := repository.db.First(&newSample, sampleID).Error; err != nil {
		return nil, nil, err
	}

	// Добавляем образец в миссию
	if err := repository.db.Create(&ds.Mission_samples{
		Id_mission: lastDraftMission.Id_mission,
		Id_sample:  newSample.Id_sample,
	}).Error; err != nil {
		return nil, nil, err
	}

	// Получаем все образцы в миссии
	var samples []ds.Samples
	err := repository.db.
		Joins("JOIN mission_samples ON missions.id_mission = mission_samples.id_mission").
		Joins("JOIN samples ON mission_samples.id_sample = samples.id_sample").
		Where("missions.id_mission = ?", lastDraftMission.Id_mission).
		Table("missions").
		Select("missions.*, samples.*").
		Find(&samples).
		Error

	if err != nil {
		return nil, nil, err
	}

	return &lastDraftMission, samples, nil

}

func (r *Repository) UpdateMissionStatusByUser(id int, newStatus string) error {
	// Проверяем, что новый статус допустим
	allowedStatus := map[string]bool{
		"Draft":                 true,
		"Deleted":               true,
		"Awaiting confirmation": true,
	}

	if !allowedStatus[newStatus] {
		return errors.New("Invalid mission status")
	}

	// Обновляем статус миссии
	err := r.db.Model(&ds.Missions{}).
		Where("Id_mission = ?", id).
		Update("mission_status", newStatus).
		Error

	return err
}

func (r *Repository) UpdateMissionStatusByModerator(id int, newStatus string) error {
	// Проверяем, что новый статус допустим
	allowedStatus := map[string]bool{
		"Completed": true,
		"Rejected":  true,
		"At work":   true,
	}

	if !allowedStatus[newStatus] {
		return errors.New("Invalid mission status")
	}

	// Обновляем статус миссии
	err := r.db.Model(&ds.Missions{}).
		Where("Id_mission = ?", id).
		Update("mission_status", newStatus).
		Error

	return err
}
func (repository *Repository) RemoveSampleFromMission(missionID, sampleID uint) (*ds.Missions, []ds.Samples, error) {
	// Проверяем, существует ли миссия с указанным ID
	var mission ds.Missions
	if err := repository.db.First(&mission, missionID).Error; err != nil {
		return nil, nil, err
	}

	// Удаляем запись об образце из таблицы mission_samples
	if err := repository.db.
		Where("id_mission = ? AND id_sample = ?", missionID, sampleID).
		Delete(&ds.Mission_samples{}).Error; err != nil {
		return nil, nil, err
	}

	// Получаем все образцы в миссии после удаления
	var samples []ds.Samples
	err := repository.db.
		Joins("JOIN mission_samples ON missions.id_mission = mission_samples.id_mission").
		Joins("JOIN samples ON mission_samples.id_sample = samples.id_sample").
		Where("missions.id_mission = ?", missionID).
		Table("missions").
		Select("missions.*, samples.*").
		Find(&samples).
		Error

	if err != nil {
		return nil, nil, err
	}

	return &mission, samples, nil
}

func (repository *Repository) RemoveSampleFromLastDraftMission(sampleID int) (*ds.Missions, []ds.Samples, error) {
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
		return nil, nil, errors.New("no draft mission found")
	}

	// Удаляем образец из миссии
	if err := repository.db.Exec("DELETE FROM mission_samples WHERE id_mission = ? AND id_sample = ?", lastDraftMission.Id_mission, sampleID).Error; err != nil {
		return nil, nil, err
	}

	// Получаем все образцы в миссии после удаления
	var samples []ds.Samples
	err := repository.db.
		Joins("JOIN mission_samples ON missions.id_mission = mission_samples.id_mission").
		Joins("JOIN samples ON mission_samples.id_sample = samples.id_sample").
		Where("missions.id_mission = ?", lastDraftMission.Id_mission).
		Table("missions").
		Select("missions.*, samples.*").
		Find(&samples).
		Error

	if err != nil {
		return nil, nil, err
	}

	return &lastDraftMission, samples, nil
}
