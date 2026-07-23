package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"service/config"
	"service/models"
)

func GetSalaryDashboard(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	unitStr := c.Query("unit")
	userStr := c.Query("user_id")

	now := time.Now()
	var startDate, endDate string

	if startDateStr != "" {
		startDate = startDateStr
	} else {
		// First day of current month
		firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = firstDay.Format("2006-01-02")
	}

	if endDateStr != "" {
		endDate = endDateStr
	} else {
		// Last day of current month
		lastDay := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location())
		endDate = lastDay.Format("2006-01-02")
	}

	// 1. Get all units
	var units []models.Unit
	if err := config.DB.Find(&units).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch units"})
		return
	}

	// 2. Build the Teach query
	// Filter whereHas('akun', function($q) { $q->where('status', 1) })
	// We join with the users table (since akun belongsTo User)
	query := config.DB.Model(&models.Teach{}).
		Joins("JOIN users ON users.id = teaches.user").
		Where("users.status = ?", 1)

	if unitStr != "" {
		unitID, err := strconv.Atoi(unitStr)
		if err == nil {
			query = query.Where("teaches.unit_id = ?", unitID)
		}
	}

	if userStr != "" {
		userID, err := strconv.Atoi(userStr)
		if err == nil {
			query = query.Where("teaches.user = ?", userID)
		}
	}

	var teaches []models.Teach
	if err := query.Find(&teaches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teachers"})
		return
	}

	// 3. For each teacher, load the filtered relations
	for i := range teaches {
		// Load unit
		var unit models.Unit
		if err := config.DB.First(&unit, teaches[i].UnitID).Error; err == nil {
			teaches[i].Unit = &unit
		}

		// Load akun (user)
		var user models.User
		if err := config.DB.First(&user, teaches[i].User).Error; err == nil {
			teaches[i].Akun = &user
			teaches[i].Name = teaches[i].Name // Hide password
		}

		// Load salaries filtered by date range
		var salaries []models.Salary
		config.DB.Where("teach_id = ? AND DATE(tanggal) >= ? AND DATE(tanggal) <= ?", teaches[i].ID, startDate, endDate).Find(&salaries)
		teaches[i].Salaries = salaries

		// Load present (StudentPresent) filtered by date range
		var presents []models.StudentPresent
		config.DB.Where("teach_id = ? AND DATE(created_at) >= ? AND DATE(created_at) <= ?", teaches[i].ID, startDate, endDate).Find(&presents)

		// For each presence record, load student, program, reg, and nested relations
		for j := range presents {
			// load student
			var student models.Student
			if err := config.DB.First(&student, presents[j].StudentID).Error; err == nil {
				presents[j].Student = &student
			}
			// load program
			var program models.Program
			if err := config.DB.First(&program, presents[j].ProgramID).Error; err == nil {
				presents[j].Program = &program
			}
			// load reg (Head)
			var head models.Head
			if err := config.DB.First(&head, presents[j].HeadID).Error; err == nil {
				// load reg.units
				var regUnit models.Unit
				if err := config.DB.First(&regUnit, head.Unit).Error; err == nil {
					head.Units = &regUnit
				}
				// load reg.product (Price)
				var price models.Price
				if err := config.DB.First(&price, head.Price).Error; err == nil {
					head.Product = &price
				}
				presents[j].Reg = &head
			}
		}
		teaches[i].Present = presents
		teaches[i].PresentCount = int64(len(presents))
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      teaches,
		"start_date": startDate,
		"end_date":   endDate,
		"units":      units,
		"unit":       unitStr,
	})
}

type GenerateSalaryRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Teachers  []struct {
		ID              uint    `json:"id"`
		Sessions        int     `json:"sessions"`
		Percentage      float64 `json:"percentage"`
		Total           float64 `json:"total"`
		JumlahPertemuan int     `json:"jumlah_pertemuan"`
	} `json:"teachers"`
}

func GenerateSalary(c *gin.Context) {
	var req GenerateSalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, t := range req.Teachers {
		tanggal, _ := time.Parse("2006-01-02", req.EndDate)

		salary := models.Salary{
			TeachID:         t.ID,
			Tanggal:         &tanggal,
			Sesi:            t.Sessions,
			Persentase:      int(t.Percentage),
			JumlahPertemuan: t.JumlahPertemuan,
			Total:           t.Total,
		}

		// GORM will create the new Salary record
		if err := config.DB.Create(&salary).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan salary untuk guru ID " + strconv.Itoa(int(t.ID))})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Salary berhasil di-generate dan tersimpan.",
	})
}
