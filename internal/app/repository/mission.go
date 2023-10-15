package repository

import (
	"MSRM/internal/app/ds"
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
	if err := r.db.Exec("UPDATE missions SET mission_status='Canceled' WHERE Id_mission= ?", id).Error; err != nil {
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
