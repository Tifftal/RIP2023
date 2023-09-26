package ds

type Samples struct {
	Id_sample        uint `gorm:"primarykey;autoIncrement"`
	Name             string
	Type             string
	Date_Sealed      string
	Sol_Sealed       int
	Rock_Type        string
	Height           string
	Current_Location string
	Image            string
	Video            string
	Sample_status    string
}
