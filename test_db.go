package main

import (
	"fmt"
	"service/config"
)

func main() {
	config.ConnectDB()
	var total int64
	q := config.DB.Table("students").
		Select("DISTINCT students.id").
		Joins("JOIN head ON head.students = students.id AND head.deleted_at IS NULL")

	q.Count(&total)

	var ids []uint
	q.Scan(&ids)
	fmt.Printf("Total: %d\nIDs: %v\n", total, ids)
}
