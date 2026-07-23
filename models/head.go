package models

type Head struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Students uint   `json:"students"`
	Program  uint   `json:"program"`
	Payment  uint   `json:"payment"`
	Price    uint   `json:"price"`
	Kelas    uint   `json:"kelas"`
	Unit     uint   `json:"unit"`
	Done     int    `json:"done"`
	Note     string `json:"note"`
	Old      int    `json:"old"`
	Number   string `json:"number"`
	Global   string `json:"global"`
	Units    *Unit    `gorm:"foreignKey:Unit" json:"units,omitempty"`
	Prices   *Price   `gorm:"foreignKey:Price" json:"prices,omitempty"`
	Product  *Price   `gorm:"foreignKey:Price" json:"product,omitempty"`
	Murid    *Student `gorm:"foreignKey:Students" json:"murid,omitempty"`
	Programs *Program `gorm:"foreignKey:Program" json:"programs,omitempty"`
}

func (Head) TableName() string {
	return "head"
}
