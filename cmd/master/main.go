package main

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/dogancankaygusuz/goprobe/internal/grpc/proto"
	"github.com/dogancankaygusuz/goprobe/pkg/database" 
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Veritabanını Başlat
	db := database.InitDB()
	log.Println("Veritabanı bağlantısı başarılı (SQLite).")

	// 2. Worker'a Bağlan
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Worker'a bağlanılamadı: %v", err)
	}
	defer conn.Close()

	client := pb.NewProbeServiceClient(conn)

	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://go.dev",
		"https://api.boredapi.com/api/activity",
	}

	for {
		log.Println("----- Taramayı Başlat -----")
		var wg sync.WaitGroup

		for _, url := range urls {
			wg.Add(1)

			go func(targetUrl string) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				resp, err := client.CheckUrl(ctx, &pb.CheckRequest{Url: targetUrl})
				
				// Veritabanı için kayıt nesnesi oluşturuyoruz
				resultToSave := database.CheckResult{
					Url: targetUrl,
				}

				if err != nil {
					log.Printf("❌ HATA [%s]: %v", targetUrl, err)
					resultToSave.Status = false
					resultToSave.ErrorMessage = err.Error()
				} else {
					log.Printf("✅ Site: %s | Kod: %d | Süre: %.0fms", 
						resp.Url, resp.StatusCode, resp.ResponseTimeMs)
					
					resultToSave.Url = resp.Url
					resultToSave.StatusCode = resp.StatusCode
					resultToSave.ResponseTimeMs = resp.ResponseTimeMs
					resultToSave.Status = resp.Status
				}

				// 3. SONUCU VERİTABANINA KAYDET (GORM ile tek satır)
				db.Create(&resultToSave)

			}(url)
		}

		wg.Wait()
		log.Println("----- Tarama Bitti ve Kaydedildi -----")
		time.Sleep(5 * time.Second)
	}
}