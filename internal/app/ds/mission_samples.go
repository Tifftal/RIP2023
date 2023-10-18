package ds

type Mission_samples struct {
	Id_mission uint `json:"Id_mission" gorm:"primaryKey"`
	Id_sample  uint `json:"Id_sample" gorm:"primaryKey"`
}
