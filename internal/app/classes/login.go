package classes

type Login struct {
	Email_address string `json:"Email_address" gorm:"column:email_address"`
	Password      string `json:"Password" gorm:"column:password"`
}
