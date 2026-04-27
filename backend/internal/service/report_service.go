package service

import (
	"backend/internal/repository"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type DashboardStats struct {
	TotalAssets  int64 `json:"totalAssets`
	ActiveAssets int64 `json:activeAssets`
	TotalValue   int64 `json:totalValue`
	TotalAudits  int64 `json: totalAudit`
}

func GetDashboardStats() (*DashboardStats, error) {
	total, _ := repository.CountAssets()
	active, _ := repository.CountActiveAssets()
	value, _ := repository.SumAssetValues()
	audits, _ := repository.CountAuditSessions()
	return &DashboardStats{
		TotalAssets: total, ActiveAssets: active,
		TotalValue: value, TotalAudits: audits,
	}, nil
}

func GenerateReport(reportType, format string) (string, error) {
	switch reportType {
	case "asset_detail":
		return generateAssetReport(format)
	case "depreciation":
		return generateDepreciationReport(format)
	default:
		return generateAssetReport(format)
	}
}

func generateAssetReport(format string) (string, error) {
	assets, err := repository.GetAllAssetsForReport()
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("/tmp/report_%s_%s.%s", "assets", uuid.New().String()[:8], format)

	switch format {
	case "csv":
		f, err := os.Create(filename)
		if err != nil {
			return "", err
		}
		defer f.Close()
		w := csv.NewWriter(f)
		w.Write([]string{"Barcode", "Name", "Price", "CurrentValue", "Status", "Location", "Date"})
		for _, a := range assets {
			loc := ""
			if a.Location != nil {
				loc = a.Location.Name
			}
			w.Write([]string{
				a.Barcode, a.AssetName,
				fmt.Sprintf("%d", a.AcquisitionPrice),
				fmt.Sprintf("%d", a.CurrentValue),
				string(a.Status), loc,
				a.AcquisitionDate.Format("2006-01-01"),
			})
		}
		w.Flush()
	case "json":
		f, err := os.Create(filename)
		if err != nil {
			return "", err
		}
		defer f.Close()
		json.NewEncoder(f).Encode(assets)
	default:
		return generateAssetReport("csv")
	}

	return filename, nil
}

func generateDepreciationReport(format string) (string, error) {
	assets, err := repository.GetAllAssetsForReport()
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("/tmp/report_depreciation_%s.csv", time.Now().Format("20060102"))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Write([]string{"Barcode", "Name", "AcquisitionPrice", "CurrentValue", "Depreciated", "UsefulLife(months)"})
	for _, a := range assets {
		dep := a.AcquisitionPrice - a.CurrentValue
		w.Write([]string{
			a.Barcode, a.AssetName,
			fmt.Sprintf("%d", a.AcquisitionPrice),
			fmt.Sprintf("%d", a.CurrentValue),
			fmt.Sprintf("%d", dep),
			fmt.Sprintf("%d", a.UsefulLifeMonths),
		})
	}
	w.Flush()
	return filename, nil
}
