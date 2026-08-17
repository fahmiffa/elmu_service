package models

import (
	"fmt"
	"strconv"
	"time"
)

type Head struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Students uint   `json:"students"`
	CreatedAt time.Time `json:"created_at"`
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

	Units    *Unit    `gorm:"foreignKey:Unit"     json:"units,omitempty"`
	Prices   *Price   `gorm:"foreignKey:Price"    json:"prices,omitempty"`
	Product  *Price   `gorm:"foreignKey:Price"    json:"product,omitempty"`
	Murid    *Student `gorm:"foreignKey:Students" json:"murid,omitempty"`
	Programs *Program `gorm:"foreignKey:Program"  json:"programs,omitempty"`
	Class    *Kelas   `gorm:"foreignKey:Kelas"    json:"class,omitempty"`

	// Relasi ke Level (banyak per head)
	Levels      []Level `gorm:"foreignKey:Head" json:"levels,omitempty"`
	LatestLevel *Level  `gorm:"-"               json:"latest_level,omitempty"`

	// Relasi ke Paid (tagihan bulanan)
	Bills []Paid `gorm:"foreignKey:Head" json:"bills,omitempty"`
}

func (h Head) GetInduk() string {
	munit := h.Number
	if len(munit) < 3 {
		munit = fmt.Sprintf("%03s", munit)
	}
	if num, err := strconv.Atoi(h.Number); err == nil {
		munit = fmt.Sprintf("%03d", num)
	}

	global := h.Global
	if len(global) < 4 {
		global = fmt.Sprintf("%04s", global)
	}
	if num, err := strconv.Atoi(h.Global); err == nil {
		global = fmt.Sprintf("%04d", num)
	}

	unitID := ""
	if h.Units != nil {
		unitID = fmt.Sprintf("%03d", h.Units.ID)
	} else if h.Unit > 0 {
		unitID = fmt.Sprintf("%03d", h.Unit)
	}

	programKode := ""
	if h.Programs != nil {
		programKode = h.Programs.Kode
	}

	return fmt.Sprintf("%s%s%s/%s", global, unitID, munit, programKode)
}

func (Head) TableName() string {
	return "head"
}
