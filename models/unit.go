package models

type Unit struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}
