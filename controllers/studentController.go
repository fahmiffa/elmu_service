package controllers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"service/config"
	"service/models"

	"github.com/gin-gonic/gin"
)

// ─── Response Structs ────────────────────────────────────────────────────────

type StudentListResponse struct {
	ID                  uint   `json:"id"`
	Name                string `json:"name"`
	NamaPanggilan       string `json:"nama_panggilan"`
	Img                 string `json:"img"`
	Jenjang             string `json:"jenjang"`
	Kelas               string `json:"kelas"`
	Place               string `json:"place"`
	Birth               string `json:"birth"`
	Gender              int    `json:"gender"`
	Genders             string `json:"genders"`
	Alamat              string `json:"alamat"`
	SekolahKelas        string `json:"sekolah_kelas"`
	AlamatSekolah       string `json:"alamat_sekolah"`
	Dream               string `json:"dream"`
	HpSiswa             string `json:"hp_siswa"`
	Agama               string `json:"agama"`
	SosmedChild         string `json:"sosmedChild"`
	SosmedOther         string `json:"sosmedOther"`
	Dad                 string `json:"dad"`
	DadJob              string `json:"dadJob"`
	Mom                 string `json:"mom"`
	MomJob              string `json:"momJob"`
	HpParent            string `json:"hp_parent"`
	Study               string `json:"study"`
	Rank                string `json:"rank"`
	PendidikanNonFormal string `json:"pendidikan_non_formal"`
	Prestasi            string `json:"prestasi"`

	UserID  uint   `json:"user_id"`
	Program string `json:"program"`
	Unit    string `json:"unit"`
	Done    int    `json:"done"`
	Status  string `json:"status"`
}

// ─── Response Structs (detail) ───────────────────────────────────────────────

type StudentUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  int    `json:"role"`
}

type StudentGradeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type StudentProgramResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Kode string `json:"kode"`
}

type StudentUnitResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type StudentKelasResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type StudentLevelResponse struct {
	ID    uint `json:"id"`
	Level int  `json:"level"`
}

type StudentRegResponse struct {
	ID      uint                    `json:"id"`
	Done    int                     `json:"done"`
	Status  string                  `json:"status"`
	Program *StudentProgramResponse `json:"program,omitempty"`
	Unit    *StudentUnitResponse    `json:"unit,omitempty"`
	Kelas   *StudentKelasResponse   `json:"kelas,omitempty"`
	Level   *StudentLevelResponse   `json:"level,omitempty"`
}

type StudentDetailResponse struct {
	ID                  uint   `json:"id"`
	Name                string `json:"name"`
	NamaPanggilan       string `json:"nama_panggilan"`
	Img                 string `json:"img"`
	Jenjang             string `json:"jenjang"`
	Kelas               string `json:"kelas"`
	Place               string `json:"place"`
	Birth               string `json:"birth"`
	Gender              int    `json:"gender"`
	Genders             string `json:"genders"`
	Alamat              string `json:"alamat"`
	SekolahKelas        string `json:"sekolah_kelas"`
	AlamatSekolah       string `json:"alamat_sekolah"`
	Dream               string `json:"dream"`
	HpSiswa             string `json:"hp_siswa"`
	Agama               string `json:"agama"`
	SosmedChild         string `json:"sosmedChild"`
	SosmedOther         string `json:"sosmedOther"`
	Dad                 string `json:"dad"`
	DadJob              string `json:"dadJob"`
	Mom                 string `json:"mom"`
	MomJob              string `json:"momJob"`
	HpParent            string `json:"hp_parent"`
	Study               string `json:"study"`
	Rank                string `json:"rank"`
	PendidikanNonFormal string `json:"pendidikan_non_formal"`
	Prestasi            string `json:"prestasi"`
	UserID              uint   `json:"user_id"`

	Akun  *StudentUserResponse  `json:"akun,omitempty"`
	Grade *StudentGradeResponse `json:"grade,omitempty"`
	Reg   []StudentRegResponse  `json:"reg"`
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func headStatusLabel(done int) string {
	switch done {
	case 0:
		return "Aktif"
	case 1:
		return "Lulus"
	case 2:
		return "Cuti"
	case 3:
		return "Keluar"
	case 4:
		return "Pindah"
	default:
		return "Tidak Diketahui"
	}
}

func genderLabel(g int) string {
	switch g {
	case 1:
		return "Laki-laki"
	case 2:
		return "Perempuan"
	default:
		return ""
	}
}

func calcTotal(p models.Paid) float64 {
	if p.Reg == nil {
		return 0
	}
	hargaBulan := 0.0
	if p.Reg.Prices != nil {
		hargaBulan = p.Reg.Prices.Harga
	}
	kit := 0.0
	if p.Reg.Programs != nil {
		kit = p.Reg.Programs.Kit
	}
	if p.First != 1 {
		return hargaBulan
	}
	created := p.CreatedAt
	if created == nil {
		return hargaBulan + kit
	}
	dayOfWeek := int(created.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}
	weekOfMonth := int(float64(created.Day()-1)/7.0) + 1
	if weekOfMonth < 1 {
		weekOfMonth = 1
	}
	if weekOfMonth > 4 {
		weekOfMonth = 4
	}
	weeksRemaining := 4 - weekOfMonth + 1
	var meetings int
	if dayOfWeek >= 1 && dayOfWeek <= 3 {
		meetings = weeksRemaining * 2
	} else {
		meetings = weeksRemaining*2 - 2
		if meetings < 0 {
			meetings = 0
		}
	}
	return float64(meetings)*(hargaBulan/8.0) + kit
}

// activeHead mengembalikan head dengan done=0 (Aktif) pertama yang ditemukan,
// atau head pertama jika tidak ada yang aktif.
func activeHead(heads []models.Head) *models.Head {
	for i := range heads {
		if heads[i].Done == 0 {
			return &heads[i]
		}
	}
	if len(heads) > 0 {
		return &heads[0]
	}
	return nil
}

// loadHeadsForStudents memuat semua head beserta Programs dan Units.
func loadHeadsForStudents(studentIDs []uint, doneFilter string) []models.Head {
	var heads []models.Head
	q := config.DB.
		Where("students IN ?", studentIDs).
		Preload("Programs").
		Preload("Units").
		Preload("Class").
		Preload("Levels")

	if doneFilter != "" {
		q = q.Where("done = ?", doneFilter)
	}
	q.Find(&heads)

	// Tentukan LatestLevel tiap head
	for i, h := range heads {
		if len(h.Levels) > 0 {
			lvls := make([]models.Level, len(h.Levels))
			copy(lvls, h.Levels)
			sort.Slice(lvls, func(a, b int) bool { return lvls[a].ID > lvls[b].ID })
			heads[i].LatestLevel = &lvls[0]
		}
	}
	return heads
}

// ─── GET /students ────────────────────────────────────────────────────────────
// Output ringkas: nama, nama_panggilan, data siswa (tanpa grade_id),
// program, unit, dan status (done) dari head aktif.
//
// Query params:
//
//	search  = nama / nama_panggilan
//	unit    = filter ID unit
//	program = filter ID program
//	done    = filter status head (kosong = semua)
//	page    = halaman (default 1)
//	limit   = per halaman (default 20, max 100)
func GetStudents(c *gin.Context) {
	search := c.Query("search")
	unitID := c.Query("unit")
	programID := c.Query("program")
	doneStr := c.Query("done")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// STEP 1: Ambil student IDs via join head
	idsQuery := config.DB.
		Table("students").
		Select("DISTINCT students.id").
		Joins("JOIN head ON head.students = students.id AND head.deleted_at IS NULL")

	if doneStr != "" {
		idsQuery = idsQuery.Where("head.done = ?", doneStr)
	}
	if unitID != "" {
		idsQuery = idsQuery.Where("head.unit = ?", unitID)
	}
	if programID != "" {
		idsQuery = idsQuery.Where("head.program = ?", programID)
	}
	if search != "" {
		s := "%" + strings.TrimSpace(search) + "%"
		idsQuery = idsQuery.Where("students.name LIKE ? OR students.nama_panggilan LIKE ?", s, s)
	}

	var total int64
	idsQuery.Count(&total)

	var studentIDs []uint
	if err := idsQuery.Offset(offset).Limit(limit).Scan(&studentIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa: " + err.Error()})
		return
	}

	if len(studentIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": []StudentListResponse{}, "total": total,
			"page": page, "limit": limit, "total_page": 0,
		})
		return
	}

	// STEP 2: Load students
	var students []models.Student
	if err := config.DB.
		Where("id IN ?", studentIDs).
		Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal load siswa: " + err.Error()})
		return
	}

	// STEP 3: Load heads (Programs + Units)
	heads := loadHeadsForStudents(studentIDs, doneStr)

	headMap := make(map[uint][]models.Head)
	for _, h := range heads {
		headMap[h.Students] = append(headMap[h.Students], h)
	}

	// STEP 4: Build ringkas response
	result := make([]StudentListResponse, 0, len(students))
	for _, s := range students {
		item := StudentListResponse{
			ID:                  s.ID,
			Name:                s.Name,
			NamaPanggilan:       s.NamaPanggilan,
			Img:                 s.Img,
			Jenjang:             s.Jenjang,
			Kelas:               s.Kelas,
			Place:               s.Place,
			Birth:               s.Birth,
			Gender:              s.Gender,
			Genders:             genderLabel(s.Gender),
			Alamat:              s.Alamat,
			SekolahKelas:        s.SekolahKelas,
			AlamatSekolah:       s.AlamatSekolah,
			Dream:               s.Dream,
			HpSiswa:             s.HpSiswa,
			Agama:               s.Agama,
			SosmedChild:         s.SosmedChild,
			SosmedOther:         s.SosmedOther,
			Dad:                 s.Dad,
			DadJob:              s.DadJob,
			Mom:                 s.Mom,
			MomJob:              s.MomJob,
			HpParent:            s.HpParent,
			Study:               s.Study,
			Rank:                s.Rank,
			PendidikanNonFormal: s.PendidikanNonFormal,
			Prestasi:            s.Prestasi,
			UserID:              s.User,
		}

		// Ambil head aktif (done=0) atau head pertama
		if h := activeHead(headMap[s.ID]); h != nil {
			item.Done = h.Done
			item.Status = headStatusLabel(h.Done)
			if h.Programs != nil {
				item.Program = h.Programs.Name
			}
			if h.Units != nil {
				item.Unit = h.Units.Name
			}
		}

		result = append(result, item)
	}

	totalPage := int(total) / limit
	if int(total)%limit > 0 {
		totalPage++
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       result,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_page": totalPage,
	})
}

// ─── GET /students/:id ────────────────────────────────────────────────────────
// Output detail: semua relasi (akun, grade, semua reg dengan program/unit/kelas/level)

func GetStudentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var student models.Student
	if err := config.DB.
		Where("id = ?", id).
		Preload("Users").
		Preload("Grade").
		First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Siswa tidak ditemukan"})
		return
	}

	heads := loadHeadsForStudents([]uint{student.ID}, "")

	sr := StudentDetailResponse{
		ID:                  student.ID,
		Name:                student.Name,
		NamaPanggilan:       student.NamaPanggilan,
		Img:                 student.Img,
		Jenjang:             student.Jenjang,
		Kelas:               student.Kelas,
		Place:               student.Place,
		Birth:               student.Birth,
		Gender:              student.Gender,
		Genders:             genderLabel(student.Gender),
		Alamat:              student.Alamat,
		SekolahKelas:        student.SekolahKelas,
		AlamatSekolah:       student.AlamatSekolah,
		Dream:               student.Dream,
		HpSiswa:             student.HpSiswa,
		Agama:               student.Agama,
		SosmedChild:         student.SosmedChild,
		SosmedOther:         student.SosmedOther,
		Dad:                 student.Dad,
		DadJob:              student.DadJob,
		Mom:                 student.Mom,
		MomJob:              student.MomJob,
		HpParent:            student.HpParent,
		Study:               student.Study,
		Rank:                student.Rank,
		PendidikanNonFormal: student.PendidikanNonFormal,
		Prestasi:            student.Prestasi,
		UserID:              student.User,
		Reg:                 []StudentRegResponse{},
	}
	if student.Users != nil {
		sr.Akun = &StudentUserResponse{
			ID:    student.Users.ID,
			Name:  student.Users.Name,
			Email: student.Users.Email,
			Role:  student.Users.Role,
		}
	}
	if student.Grade != nil {
		sr.Grade = &StudentGradeResponse{ID: student.Grade.ID, Name: student.Grade.Name}
	}

	for _, h := range heads {
		reg := StudentRegResponse{
			ID:     h.ID,
			Done:   h.Done,
			Status: headStatusLabel(h.Done),
		}
		if h.Programs != nil {
			reg.Program = &StudentProgramResponse{ID: h.Programs.ID, Name: h.Programs.Name, Kode: h.Programs.Kode}
		}
		if h.Units != nil {
			reg.Unit = &StudentUnitResponse{ID: h.Units.ID, Name: h.Units.Name}
		}
		if h.Class != nil {
			reg.Kelas = &StudentKelasResponse{ID: h.Class.ID, Name: h.Class.Name}
		}
		if h.LatestLevel != nil {
			reg.Level = &StudentLevelResponse{ID: h.LatestLevel.ID, Level: h.LatestLevel.Level}
		}
		sr.Reg = append(sr.Reg, reg)
	}

	c.JSON(http.StatusOK, gin.H{"data": sr})
}
