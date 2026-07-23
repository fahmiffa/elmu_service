package models

type Price struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	Kelas   uint    `json:"kelas"`
	Product uint    `json:"product"`
	Harga   float64 `json:"harga"`
}
