package models

type Tabsensor struct {
	Sourceid  int    `gorm:"column:sourceid" json:"sourceid"`
	Latitude  string `gorm:"column:latitude" json:"latitude"`
	Longitude string `gorm:"column:longitude" json:"longitude"`
}

func (Tabsensor) TableName() string {
	return "tabsensor"
}
