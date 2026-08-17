package models

import "time"

type StudentPresent struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	StudentID       uint           `json:"student_id"`
	UnitSchedulesID uint           `json:"unit_schedules_id"`
	TeachID         uint           `json:"teach_id"`
	Hal             string         `json:"hal"`
	Materi          string         `json:"materi"`
	Keterangan      string         `json:"keterangan"`
	HeadID          uint           `json:"head_id"`
	Meet    		 int		   `json:"meet"`
	Present 		 int		   `json:"present"`
	ProgramID       uint           `json:"program_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Student         *Student       `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Program         *Program       `gorm:"foreignKey:ProgramID" json:"program,omitempty"`
	Reg             *Head          `gorm:"foreignKey:HeadID" json:"reg,omitempty"`
}

func (StudentPresent) TableName() string {
	return "student_presents"
}
