package ds

type Users struct {
	Id_user        uint   `json:"Id_user" gorm:"primarykey;autoIncrement"`
	Name           string `json:"Name" gorm:"column:name"`
	Phone_number   string `json:"Phone_number" gorm:"column:phone_number"`
	Email_address  string `json:"Email_address" gorm:"column:email_address"`
	Password       string `json:"Password" gorm:"column:password"`
	RepeatPassword string `json:"RepeatPassword" gorm:"column:repeat_password"`
	Role           string `json:"Role" gorm:"column:role"`
	User_status    string `json:"User_status" gorm:"column:user_status"`
}
