package models

import "time"

type Student struct {
	ID                  uint   `gorm:"primaryKey" json:"id"`
	User                uint   `json:"user"`
	Name                string `json:"name"`
	NamaPanggilan       string `json:"nama_panggilan"`
	Img                 string `json:"img"`
	Jenjang             string `json:"jenjang"`
	Kelas               string `json:"kelas"`
	Place               string `json:"place"`
	Birth               string `json:"birth"`
	Gender              int    `json:"gender"`
	GradeID             uint   `json:"grade_id"`
	Alamat              string `json:"alamat"`
	SekolahKelas        string `json:"sekolah_kelas"`
	AlamatSekolah       string `json:"alamat_sekolah"`
	Dream               string `json:"dream"`
	HpSiswa             string `json:"hp_siswa"`
	Agama               string `json:"agama"`
	SosmedChild         string `gorm:"column:sosmedChild" json:"sosmedChild"`
	SosmedOther         string `gorm:"column:sosmedOther" json:"sosmedOther"`
	Dad                 string `json:"dad"`
	DadJob              string `gorm:"column:dadJob" json:"dadJob"`
	Mom                 string `json:"mom"`
	MomJob              string `gorm:"column:momJob" json:"momJob"`
	HpParent            string `json:"hp_parent"`
	Study               string `json:"study"`
	Rank                string `json:"rank"`
	PendidikanNonFormal string `json:"pendidikan_non_formal"`
	Prestasi            string `json:"prestasi"`
	CreatedAt           time.Time `json:"created_at"`

	// Relasi
	Users *User  `gorm:"foreignKey:User" json:"users,omitempty"`
	Grade *Grade `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
	Reg   []Head `gorm:"foreignKey:Students" json:"reg,omitempty"`
}

func (Student) TableName() string {
	return "students"
}
