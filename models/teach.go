package models

type Teach struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	User         uint             `json:"user"`
	Name         string           `json:"name"`
	UnitID       uint             `json:"unit_id"`
	Hp           string           `json:"hp"`
	Img          string           `json:"img"`
	Study        string           `json:"study"`
	Addr         string           `json:"addr"`
	Birth        string           `json:"birth"`
	Profit       int              `json:"profit"`
	Unit         *Unit            `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	Akun         *User            `gorm:"foreignKey:User" json:"akun,omitempty"`
	Salaries     []Salary         `gorm:"foreignKey:TeachID" json:"salaries,omitempty"`
	Present      []StudentPresent `gorm:"foreignKey:TeachID" json:"present,omitempty"`
	PresentCount int64            `gorm:"-" json:"present_count"`
}
