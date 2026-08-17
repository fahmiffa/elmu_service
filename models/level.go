package models

import "time"

type Level struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Head      uint       `json:"head"`
	TeachUser uint       `gorm:"column:teach_user" json:"teach_user"`
	Level     string     `json:"level"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (Level) TableName() string {
	return "levels"
}
