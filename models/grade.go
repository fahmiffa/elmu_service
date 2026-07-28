package models

type Grade struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

func (Grade) TableName() string {
	return "grades"
}
