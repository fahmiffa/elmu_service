package models

import "time"

type Order struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Head      uint       `json:"head"`
	Student   uint       `json:"student"`
	Price     uint       `json:"price"`
	Status    int        `json:"status"`
	Via       string     `json:"via"`
	Mid       string     `json:"mid"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	Reg     *Head  `gorm:"foreignKey:Head" json:"reg,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}
