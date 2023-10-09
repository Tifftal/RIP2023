package repository

import (
	"MSRM/internal/app/ds"
)

func (repository *Repository) GetAllMissiions() ([]ds.Missions, error) {
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
