package models

type Student struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	User          uint   `json:"user"`
	Name          string `json:"name"`
	NamaPanggilan string `json:"nama_panggilan"`
	Birth         string `json:"birth"`
	Gender        int    `json:"gender"`
	GradeID       uint   `json:"grade_id"`
}
