package ds

import (
	"time"
)

type Missions struct {
	Id_mission      uint `json:"Id_mission" gorm:"primarykey;autoIncrement"`
	User_id         int
	Moderator_id    int
	Name            string     `json:"Name" gorm:"column:name"`
	Mission_status  string     `json:"Mission_status" gorm:"column:mission_status"`
	Creation_date   *time.Time `json:"Creation_date"`
	Formation_date  *time.Time `json:"Formation_date"`
	Completion_date *time.Time `json:"Completion_date"`
	Samples         []Samples  `gorm:"many2many:mission_samples;joinForeignKey:Id_mission;joinReferences:Id_mission" json:"samples"`
}
