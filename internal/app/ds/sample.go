package ds

import "time"

type Samples struct {
	Id_sample        uint      `json:"Id_sample" gorm:"primarykey;autoIncrement"`
	Name             string    `json:"Name" gorm:"column:name"`
	Type             string    `json:"Type" gorm:"column:type"`
	Date_Sealed      time.Time `json:"Date_Sealed"`
	Sol_Sealed       int       `json:"Sol_Sealed" gorm:"column:sol_sealed"`
	Rock_Type        string    `json:"Rock_Type" gorm:"column:rock_type"`
	Height           string    `json:"Height" gorm:"column:height"`
	Current_Location string    `json:"Current_Location" gorm:"column:current_location"`
	Image            string    `json:"Image" gorm:"column:image"`
	Video            string    `json:"Video" gorm:"column:video"`
	Sample_status    string    `json:"Sample_status" gorm:"column:sample_status"`
}
