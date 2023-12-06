package classes

type UpdateMissionStatus struct {
	Mission_status string `json:"Mission_status" gorm:"column:mission_status"`
}
