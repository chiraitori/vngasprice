package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"vn-petro-price/scraper"
)

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Get prices from both Petrolimex and PVOil
	r.GET("/api/v1/prices", func(c *gin.Context) {
		var wg sync.WaitGroup
		var petrolimexData, pvoilData scraper.PriceData

		wg.Add(2)

		go func() {
			defer wg.Done()
			petrolimexData = scraper.ScrapePetrolimex()
		}()

		go func() {
			defer wg.Done()
			pvoilData = scraper.ScrapePVOil()
		}()

		wg.Wait()

		response := scraper.PriceResponse{
			Petrolimex: petrolimexData,
			PVOil:      pvoilData,
		}

		c.JSON(http.StatusOK, response)
	})

	// Get prices from Petrolimex only
	r.GET("/api/v1/prices/petrolimex", func(c *gin.Context) {
		data := scraper.ScrapePetrolimex()
		c.JSON(http.StatusOK, data)
	})

	// Get prices from PVOil only
	r.GET("/api/v1/prices/pvoil", func(c *gin.Context) {
		data := scraper.ScrapePVOil()
		c.JSON(http.StatusOK, data)
	})

	log.Println("🚀 VN Petro Price API starting on :6320")
	log.Println("📍 GET /api/v1/prices          - All prices (Petrolimex + PVOil)")
	log.Println("📍 GET /api/v1/prices/petrolimex - Petrolimex only")
	log.Println("📍 GET /api/v1/prices/pvoil      - PVOil only")
	log.Println("📍 GET /health                   - Health check")

	if err := r.Run(":6320"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
