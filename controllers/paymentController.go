package controllers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"service/config"
	"service/models"

	"github.com/gin-gonic/gin"
)

func GetMonthlyPayments(c *gin.Context) {
	tab := c.DefaultQuery("tab", "tagihan")
	bulanStr := c.Query("bulan")
	tahunStr := c.Query("tahun")
	unitID := c.Query("unit")
	programID := c.Query("program")
	search := c.Query("search")

	now := time.Now()
	if bulanStr == "" {
		bulanStr = now.Format("01")
	}
	if tahunStr == "" {
		tahunStr = strconv.Itoa(now.Year())
	}

	// STEP 1: Get paid IDs (only real DB columns - total is computed, not stored)
	idsQuery := config.DB.
		Table("paids").
		Select("paids.id").
		Joins("LEFT JOIN head ON head.id = paids.head").
		Joins("LEFT JOIN students ON students.id = head.students").
		Where("paids.bulan = ? AND paids.tahun = ?", bulanStr, tahunStr)

	if tab == "riwayat" {
		idsQuery = idsQuery.Where("paids.status = ?", 1)
	} else {
		idsQuery = idsQuery.Where("paids.status != ?", 1)
	}

	if unitID != "" {
		idsQuery = idsQuery.Where("head.unit = ?", unitID)
	}

	if programID != "" {
		idsQuery = idsQuery.Where("head.program = ?", programID)
	}

	if search != "" {
		idsQuery = idsQuery.Where("students.name LIKE ? OR students.nama_panggilan LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var ids []uint
	if err := idsQuery.Find(&ids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments: " + err.Error()})
		return
	}

	// STEP 2: Load full Paid records (only real columns, excluding 'total')
	var paids []models.Paid
	if len(ids) > 0 {
		if err := config.DB.
			Select("id, head, bulan, tahun, status, first, via, mid, created_at, updated_at").
			Preload("Reg.Murid").
			Preload("Reg.Programs").
			Preload("Reg.Units").
			Preload("Reg.Prices").
			Where("id IN ?", ids).
			Order("paids.id ASC").
			Find(&paids).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments: " + err.Error()})
			return
		}
	}

	// STEP 3: Load filters
	var units []models.Unit
	config.DB.Find(&units)

	var programs []models.Program
	config.DB.Find(&programs)

	// STEP 4: Build enriched response with computed 'total'
	type PaymentItem struct {
		ID        uint                   `json:"id"`
		Head      uint                   `json:"head"`
		Bulan     string                 `json:"bulan"`
		Tahun     string                 `json:"tahun"`
		Status    int                    `json:"status"`
		First     int                    `json:"first"`
		Via       string                 `json:"via"`
		Mid       string                 `json:"mid"`
		Total     float64                `json:"total"`
		Tempo     string                 `json:"tempo"`
		Tipe      string                 `json:"tipe"`
		CreatedAt *time.Time             `json:"created_at"`
		UpdatedAt *time.Time             `json:"updated_at"`
		Reg       map[string]interface{} `json:"reg,omitempty"`
	}

	var items []PaymentItem
	for _, p := range paids {
		total := calculateTotal(p)

		reg := map[string]interface{}{
			"id": p.Head,
		}
		if p.Reg != nil {
			if p.Reg.Murid != nil {
				reg["murid"] = map[string]interface{}{
					"id":             p.Reg.Murid.ID,
					"name":           p.Reg.Murid.Name,
					"nama_panggilan": p.Reg.Murid.NamaPanggilan,
				}
			}
			if p.Reg.Programs != nil {
				reg["programs"] = map[string]interface{}{
					"id":   p.Reg.Programs.ID,
					"name": p.Reg.Programs.Name,
				}
			}
			if p.Reg.Units != nil {
				reg["units"] = map[string]interface{}{
					"id":   p.Reg.Units.ID,
					"name": p.Reg.Units.Name,
				}
			}
		}

		item := PaymentItem{
			ID:        p.ID,
			Head:      p.Head,
			Bulan:     p.Bulan,
			Tahun:     p.Tahun,
			Status:    p.Status,
			First:     p.First,
			Via:       p.Via,
			Mid:       p.Mid,
			Total:     total,
			Tempo:     p.Bulan + "/" + p.Tahun,
			Tipe:      "Bulanan",
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			Reg:       reg,
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"items":    items,
		"units":    units,
		"programs": programs,
		"bulan":    bulanStr,
		"tahun":    tahunStr,
		"tab":      tab,
	})
}

// calculateTotal replicates Laravel's Paid::getTotalAttribute()
func calculateTotal(p models.Paid) float64 {
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

	// Not first registration -> full monthly price
	if p.First != 1 {
		return hargaBulan
	}

	// First registration: prorated based on created_at
	created := p.CreatedAt
	if created == nil {
		return hargaBulan + kit
	}

	// Day of week (1=Monday, 7=Sunday)
	day := created.Weekday()
	dayOfWeek := int(day)
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	// Approximate week of month
	weekOfMonth := int(math.Ceil(float64(created.Day()) / 7.0))
	if weekOfMonth < 1 {
		weekOfMonth = 1
	}
	if weekOfMonth > 4 {
		weekOfMonth = 4
	}

	weeksRemaining := 4 - weekOfMonth + 1

	var meetings int
	if dayOfWeek >= 1 && dayOfWeek <= 3 {
		// Monday-Wednesday: full remaining weeks
		meetings = weeksRemaining * 2
	} else {
		// Thursday-Sunday: minus 2 for current partial week
		meetings = weeksRemaining*2 - 2
		if meetings < 0 {
			meetings = 0
		}
	}

	pricePerMeeting := hargaBulan / 8.0
	total := float64(meetings) * pricePerMeeting

	// Add kit cost for first registration
	total += kit

	return math.Round(total)
}
