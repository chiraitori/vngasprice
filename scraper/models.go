package scraper

// FuelPrice represents the price of a specific fuel type
type FuelPrice struct {
	Name       string `json:"name"`
	PriceZone1 string `json:"price_zone1"`
	PriceZone2 string `json:"price_zone2,omitempty"`
}

// PriceResponse holds the scraped prices from both companies
type PriceResponse struct {
	Petrolimex PriceData `json:"petrolimex"`
	PVOil      PriceData `json:"pvoil"`
}

// PriceData holds fuel prices and metadata for a single company
type PriceData struct {
	Company   string      `json:"company"`
	UpdatedAt string      `json:"updated_at,omitempty"`
	Prices    []FuelPrice `json:"prices"`
	Error     string      `json:"error,omitempty"`
}
