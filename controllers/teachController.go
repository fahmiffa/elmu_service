package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"service/config"
	"service/models"

	"github.com/gin-gonic/gin"
)

// ─── Response Structs ────────────────────────────────────────────────────────

type TeachListResponse struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Hp      string  `json:"hp"`
	Img     string  `json:"img"`
	Study   string  `json:"study"`
	Addr    string  `json:"addr"`
	Birth   string  `json:"birth"`
	Profit  float64 `json:"profit"`
	UserID  uint    `json:"user_id"`
	Akun    string  `json:"akun"`
	Role    int     `json:"role"`
	Unit    string  `json:"unit"`
}

type TeachDetailResponse struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Hp      string  `json:"hp"`
	Img     string  `json:"img"`
	Study   string  `json:"study"`
	Addr    string  `json:"addr"`
	Birth   string  `json:"birth"`
	Profit  float64 `json:"profit"`
	UserID  uint    `json:"user_id"`
	Akun    *models.User `json:"akun,omitempty"`
	Unit    *models.Unit `json:"unit,omitempty"`
}

// ─── GET /teachers ────────────────────────────────────────────────────────────
func GetTeachers(c *gin.Context) {
	search := c.Query("search")
	unitID := c.Query("unit")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := config.DB.Model(&models.Teach{})

	if unitID != "" {
		query = query.Where("unit_id = ?", unitID)
	}
	if search != "" {
		s := "%" + strings.TrimSpace(search) + "%"
		query = query.Where("name LIKE ?", s)
	}

	var total int64
	query.Count(&total)

	var teaches []models.Teach
	if err := query.Preload("Akun").Preload("Unit").Offset(offset).Limit(limit).Find(&teaches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data guru: " + err.Error()})
		return
	}

	result := make([]TeachListResponse, 0, len(teaches))
	for _, t := range teaches {
		item := TeachListResponse{
			ID:     t.ID,
			Name:   t.Name,
			Hp:     t.Hp,
			Img:    t.Img,
			Study:  t.Study,
			Addr:   t.Addr,
			Birth:  t.Birth,
			Profit: float64(t.Profit),
			UserID: t.User,
		}

		if t.Akun != nil {
			item.Akun = t.Akun.Name
			item.Role = t.Akun.Role
		}
		if t.Unit != nil {
			item.Unit = t.Unit.Name
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

// ─── GET /teachers/:id ────────────────────────────────────────────────────────
func GetTeacherByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var teach models.Teach
	if err := config.DB.
		Where("id = ?", id).
		Preload("Akun").
		Preload("Unit").
		First(&teach).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Guru tidak ditemukan"})
		return
	}

	sr := TeachDetailResponse{
		ID:     teach.ID,
		Name:   teach.Name,
		Hp:     teach.Hp,
		Img:    teach.Img,
		Study:  teach.Study,
		Addr:   teach.Addr,
		Birth:  teach.Birth,
		Profit: float64(teach.Profit),
		UserID: teach.User,
		Akun:   teach.Akun,
		Unit:   teach.Unit,
	}

	c.JSON(http.StatusOK, gin.H{"data": sr})
}
