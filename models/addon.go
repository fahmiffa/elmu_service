package models

type Addon struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

func (Addon) TableName() string {
	return "addons"
}
