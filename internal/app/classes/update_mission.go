package classes

import "time"

type UpdateMission struct {
	Name            string     `json:"Name" gorm:"column:name"`
	Formation_date  *time.Time `json:"Formation_date"`
	Completion_date *time.Time `json:"Completion_date"`
}
