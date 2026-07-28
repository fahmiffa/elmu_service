package models

type Kelas struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

func (Kelas) TableName() string {
	return "kelas"
}
