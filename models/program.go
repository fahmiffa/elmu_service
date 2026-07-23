package models

type Program struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Name   string  `json:"name"`
	Kode   string  `json:"kode"`
	Des    string  `json:"des"`
	Level  int     `json:"level"`
	Kit    float64 `json:"kit"`
	KitDes string  `gorm:"column:kit_des" json:"kit_des"`
	Extend *bool   `gorm:"default:null" json:"extend"`
}
