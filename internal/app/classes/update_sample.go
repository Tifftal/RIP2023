package classes

type Sample_Update struct {
	Name             string `json:"Name" gorm:"column:name"`
	Type             string `json:"Type" gorm:"column:type"`
	Rock_Type        string `json:"Rock_Type" gorm:"column:rock_type"`
	Current_Location string `json:"Current_Location" gorm:"column:current_location"`
	Sample_status    string `json:"Sample_status" gorm:"column:sample_status"`
}
