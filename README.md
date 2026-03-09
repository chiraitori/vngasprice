# vngasprice

A Go REST API that scrapes and serves Vietnam fuel prices from **Petrolimex** and **PVOil**.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/prices` | Combined prices from both Petrolimex and PVOil |
| `GET` | `/api/v1/prices/petrolimex` | Petrolimex prices only (Zone 1 + Zone 2) |
| `GET` | `/api/v1/prices/pvoil` | PVOil prices only (Zone 1) |
| `GET` | `/health` | Health check |

## How to Run

```bash
go run main.go
# API available at http://localhost:6320
```
