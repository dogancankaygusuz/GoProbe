package database

import (
	"log"

	"github.com/glebarez/sqlite" // Pure Go SQLite sürücüsü
	"gorm.io/gorm"
)

// CheckResult: Veritabanına kaydedilecek tablo yapısı (Model)
type CheckResult struct {
	gorm.Model     // ID, CreatedAt, UpdatedAt, DeletedAt alanlarını otomatik ekler
	Url            string
	StatusCode     int32
	ResponseTimeMs float64
	Status         bool   // Site ayakta mı?
	ErrorMessage   string // Hata mesajı (varsa)
}

// InitDB: Veritabanını başlatır ve tabloları oluşturur
func InitDB() *gorm.DB {
	// goprobe.db adında bir dosya oluşturacak
	db, err := gorm.Open(sqlite.Open("goprobe.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	// AutoMigrate: CheckResult struct'ına bakarak veritabanında tabloyu otomatik oluşturur
	err = db.AutoMigrate(&CheckResult{})
	if err != nil {
		log.Fatalf("Tablo oluşturulamadı: %v", err)
	}

	return db
}
