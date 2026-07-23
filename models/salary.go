package models

import "time"

type Salary struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	TeachID         uint       `json:"teach_id"`
	Nominal         float64    `json:"nominal"`
	Status          int        `json:"status"`
	Tanggal         *time.Time `gorm:"type:date" json:"tanggal"`
	Sesi            int        `json:"sesi"`
	Persentase      int        `json:"persentase"`
	Total           float64    `json:"total"`
	JumlahPertemuan int        `json:"jumlah_pertemuan"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
