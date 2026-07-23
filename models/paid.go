package models

import "time"

type Paid struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Head      uint       `json:"head"`
	Bulan     string     `json:"bulan"`
	Tahun     string     `json:"tahun"`
	Status    int        `json:"status"`
	First     int        `json:"first"`
	Via       string     `json:"via"`
	Mid       string     `json:"mid"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	// NOTE: 'total' tidak ada di tabel database.
	// Ini adalah computed attribute dari Laravel (getTotalAttribute),
	// dihitung dari harga prices dan kit programs.
	// Kita akan hitung manual di controller.

	Reg *Head `gorm:"foreignKey:Head" json:"reg,omitempty"`
}

func (Paid) TableName() string {
	return "paids"
}
