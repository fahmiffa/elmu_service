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
	studentID := c.Query("student")

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

	if studentID != "" {
		idsQuery = idsQuery.Where("students.id = ?", studentID)
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

// VerifyPayment updates the status of a monthly payment to 1 (Paid)
func VerifyPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var paid models.Paid
	if err := config.DB.First(&paid, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tagihan tidak ditemukan"})
		return
	}

	if paid.Status == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tagihan ini sudah dibayar"})
		return
	}

	paid.Status = 1
	if err := config.DB.Save(&paid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pembayaran berhasil diverifikasi",
		"data":    paid,
	})
}

// GetServicePayments fetches 'Layanan' payments
func GetServicePayments(c *gin.Context) {
	tab := c.DefaultQuery("tab", "tagihan")
	search := c.Query("search")
	unitID := c.Query("unit")
	programID := c.Query("program")
	studentID := c.Query("student")

	// Custom struct to hold the joined result
	type ServiceItem struct {
		ID        uint       `json:"id"`
		Status    int        `json:"status"`
		Total     float64    `json:"total"`
		Via       string     `json:"via"`
		CreatedAt *time.Time `json:"created_at"`
		Reg       map[string]interface{} `json:"reg"`
	}

	query := config.DB.Table("orders").
		Select(`
			orders.id, orders.status, orders.via, orders.created_at,
			prices.harga as total,
			addons.name as layanan_name,
			students.id as murid_id, students.name as murid_name, students.nama_panggilan,
			units.name as unit_name,
			programs.name as program_name
		`).
		Joins("LEFT JOIN prices ON prices.id = orders.price").
		Joins("LEFT JOIN addons ON addons.id = prices.product").
		Joins("LEFT JOIN head ON head.id = orders.head").
		Joins("LEFT JOIN students ON students.id = head.students").
		Joins("LEFT JOIN units ON units.id = head.unit").
		Joins("LEFT JOIN programs ON programs.id = head.program")

	if tab == "riwayat" {
		query = query.Where("orders.status = ?", 1)
	} else {
		query = query.Where("orders.status != ?", 1)
	}

	if search != "" {
		query = query.Where("students.name LIKE ? OR students.nama_panggilan LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if unitID != "" {
		query = query.Where("head.unit = ?", unitID)
	}
	if programID != "" {
		query = query.Where("head.program = ?", programID)
	}
	if studentID != "" {
		query = query.Where("students.id = ?", studentID)
	}

	rows, err := query.Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch service payments"})
		return
	}
	defer rows.Close()

	var items []ServiceItem
	for rows.Next() {
		var id uint
		var status int
		var total float64
		var via string
		var createdAt *time.Time
		var layananName, muridIdStr, muridName, namaPanggilan, unitName, programName string

		config.DB.ScanRows(rows, &id) // Simplified scan via manual struct or variables
		rows.Scan(&id, &status, &via, &createdAt, &total, &layananName, &muridIdStr, &muridName, &namaPanggilan, &unitName, &programName)

		item := ServiceItem{
			ID: id, Status: status, Total: total, Via: via, CreatedAt: createdAt,
			Reg: map[string]interface{}{
				"layanan_name": layananName,
				"murid": map[string]interface{}{
					"name": muridName,
					"nama_panggilan": namaPanggilan,
				},
				"units": map[string]interface{}{"name": unitName},
				"programs": map[string]interface{}{"name": programName},
			},
		}
		items = append(items, item)
	}

	// Fetch filters
	var units []models.Unit
	config.DB.Find(&units)

	var programs []models.Program
	config.DB.Find(&programs)

	c.JSON(http.StatusOK, gin.H{
		"items":    items,
		"units":    units,
		"programs": programs,
		"tab":      tab,
	})
}

// VerifyServicePayment updates order status to 1
func VerifyServicePayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tagihan layanan tidak ditemukan"})
		return
	}

	if order.Status == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tagihan layanan ini sudah dibayar"})
		return
	}

	order.Status = 1
	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi pembayaran layanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran layanan berhasil diverifikasi"})
}
