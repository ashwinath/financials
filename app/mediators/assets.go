package mediator

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
	"github.com/ashwinath/financials/api/utils"
)

// AssetMediator is the mediator that calculates and stores the assets.
type AssetMediator struct {
	assetService  *service.AssetService
	assetsCSVFile string
}

// NewAssetMediator creates a new AssetMediator
func NewAssetMediator(
	assetService *service.AssetService,
	assetsCSVFile string,
) *AssetMediator {
	return &AssetMediator{
		assetService:  assetService,
		assetsCSVFile: assetsCSVFile,
	}
}

func (m *AssetMediator) parseCSV() ([]*models.Asset, error) {
	records, err := utils.ReadCSV(m.assetsCSVFile)
	if err != nil {
		return nil, err
	}

	var assets []*models.Asset
	headers := records[0]
	for recordNum := 1; recordNum < len(records); recordNum++ {
		asset := &models.Asset{}
		for i, value := range records[recordNum] {
			switch headers[i] {
			case "date":
				layout := "2006-01-02T15:04:05.000Z"
				str := fmt.Sprintf("%sT08:00:00.000Z", value)
				t, err := time.Parse(layout, str)
				if err != nil {
					return nil, err
				}
				asset.TransactionDate = t
			case "type":
				asset.Type = value
			case "amount":
				if v, err := strconv.ParseFloat(value, 64); err == nil {
					asset.Amount = v
				} else {
					return nil, err
				}
			}
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

// ProcessAssets reads the csvs and stores them
func (m *AssetMediator) ProcessAssets() error {
	log.Printf("Updating assets.")
	assets, err := m.parseCSV()
	if err != nil {
		return err
	}

	err = m.assetService.TruncateTable()
	if err != nil {
		return err
	}

	err = m.assetService.BulkAdd(assets)
	if err != nil {
		return err
	}

	return nil
}
