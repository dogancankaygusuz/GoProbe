package main

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/dogancankaygusuz/goprobe/internal/grpc/proto"
	"github.com/dogancankaygusuz/goprobe/pkg/config"
	"github.com/dogancankaygusuz/goprobe/pkg/database"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Konfigürasyonu Yükle
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Ayar dosyası (config.json) okunamadı: %v", err)
	}
	log.Printf("Konfigürasyon yüklendi. İzlenecek site sayısı: %d", len(cfg.Targets))

	// Veritabanını Başlat (SQLite)
	db := database.InitDB()
	log.Println("Veritabanı bağlantısı başarılı (SQLite).")

	// Worker'a Bağlan
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Worker'a bağlanılamadı: %v", err)
	}
	defer conn.Close()

	client := pb.NewProbeServiceClient(conn)

	// Sonsuz Döngüde Tarama
	for {
		log.Println("----- Taramayı Başlat -----")
		startTotal := time.Now()
		var wg sync.WaitGroup

		// Config'den gelen URL listesini dönüyoruz
		for _, url := range cfg.Targets {
			wg.Add(1)

			// Her site için bir Goroutine
			go func(targetUrl string) {
				defer wg.Done()

				// Config'den gelen timeout süresini kullanıyoruz
				timeoutDuration := time.Duration(cfg.Timeout) * time.Second
				ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
				defer cancel()

				// gRPC ile Worker'a sor
				resp, err := client.CheckUrl(ctx, &pb.CheckRequest{Url: targetUrl})

				// Veritabanı kayıt nesnesi
				resultToSave := database.CheckResult{
					Url: targetUrl,
				}

				if err != nil {
					// Hata durumu
					log.Printf("❌ HATA [%s]: %v", targetUrl, err)
					resultToSave.Status = false
					resultToSave.ErrorMessage = err.Error()
				} else {
					// Başarılı durum
					statusIcon := "✅"
					if !resp.Status {
						statusIcon = "🔻"
					}
					log.Printf("%s Site: %s | Kod: %d | Süre: %.0fms",
						statusIcon, resp.Url, resp.StatusCode, resp.ResponseTimeMs)

					resultToSave.Url = resp.Url
					resultToSave.StatusCode = resp.StatusCode
					resultToSave.ResponseTimeMs = resp.ResponseTimeMs
					resultToSave.Status = resp.Status
				}

				// Sonucu Veritabanına Kaydet
				db.Create(&resultToSave)

			}(url)
		}
		wg.Wait()
		totalDuration := time.Since(startTotal)
		log.Printf("----- Tarama Bitti (Toplam Süre: %v) -----", totalDuration)
		time.Sleep(5 * time.Second)
	}
}
