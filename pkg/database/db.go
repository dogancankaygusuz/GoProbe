package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Veritabanına kaydedilecek tablo yapısı
type CheckResult struct {
	gorm.Model     // ID, CreatedAt, UpdatedAt, DeletedAt alanlarını otomatik ekler
	Url            string
	StatusCode     int32
	ResponseTimeMs float64
	Status         bool
	ErrorMessage   string
}

// Veritabanını başlatır ve tabloları oluşturur
func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("goprobe.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	err = db.AutoMigrate(&CheckResult{})
	if err != nil {
		log.Fatalf("Tablo oluşturulamadı: %v", err)
	}
	return db
}
