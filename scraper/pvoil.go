package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const pvoilURL = "https://www.pvoil.com.vn/"

// ScrapePVOil fetches the latest retail fuel prices from PVOil.
// PVOil renders prices server-side on their homepage inside
// section.home-gas-price > div.gasoline-prices-list > a.gasoline-price-item
func ScrapePVOil() PriceData {
	data := PriceData{
		Company: "PVOil",
	}

	req, err := http.NewRequest("GET", pvoilURL, nil)
	if err != nil {
		data.Error = fmt.Sprintf("failed to create request: %v", err)
		return data
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en;q=0.8")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		data.Error = fmt.Sprintf("failed to fetch: %v", err)
		return data
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		data.Error = fmt.Sprintf("status code: %d", res.StatusCode)
		return data
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		data.Error = fmt.Sprintf("failed to parse HTML: %v", err)
		return data
	}

	// Extract update timestamp from h4.sub-box-title
	doc.Find("h4.sub-box-title").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "Giá điều chỉnh") || strings.Contains(text, "điều chỉnh") {
			data.UpdatedAt = text
		}
	})

	// Extract fuel prices from a.gasoline-price-item elements
	doc.Find("a.gasoline-price-item").Each(func(i int, s *goquery.Selection) {
		// Fuel name is in h3.title
		nameEl := s.Find("h3.title")
		if nameEl.Length() == 0 {
			return
		}
		name := strings.TrimSpace(nameEl.Text())

		// Price is in span.count
		priceEl := s.Find("span.count")
		if priceEl.Length() == 0 {
			return
		}
		priceText := strings.TrimSpace(priceEl.Text())
		priceText = strings.ReplaceAll(priceText, "đ", "")
		priceText = strings.TrimSpace(priceText)

		if name != "" && priceText != "" {
			data.Prices = append(data.Prices, FuelPrice{
				Name:       name,
				PriceZone1: priceText,
			})
		}
	})

	if len(data.Prices) == 0 {
		data.Error = "no prices extracted from page"
	}

	return data
}
