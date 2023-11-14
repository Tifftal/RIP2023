package ds

type Mission_samples struct {
	Mission_id uint `json:"Id_mission" gorm:"primarykey"`
	Sample_id  uint `json:"Id_sample" gorm:"primarykey"`
}
