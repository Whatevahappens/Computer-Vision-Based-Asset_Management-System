package service

import (
	"backend/internal/dto"
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

func CalculateDepreciation(assetID, method, userID string) (*dto.DepreciationResult, error) {
	asset, err := repository.FindAssetByID(assetID)
	if err != nil {
		return nil, fmt.Errorf("asset not found")
	}

	var monthly float64

	switch model.DepreciationMethod(method) {
	case model.StraightLine:
		if asset.UsefulLifeMonths <= 0 {
			return nil, fmt.Errorf("useful life must be > 0")
		}
		monthly = float64(asset.AcquisitionPrice-0) / float64(asset.UsefulLifeMonths)

	case model.DecliningBalance:
		years := float64(asset.UsefulLifeMonths) / 12.0
		if years <= 0 {
			return nil, fmt.Errorf("useful life must be > 0")
		}
		rate := 2.0 / years
		monthly = float64(asset.CurrentValue) * rate / 12.0

	default:
		return nil, fmt.Errorf("unsupported depreciation method: %s", method)
	}

	monthly = math.Round(monthly*100) / 100
	newValue := asset.CurrentValue - int(monthly)
	if newValue < 0 {
		newValue = 0
	}
	asset.CurrentValue = newValue
	repository.UpdateAsset(asset)

	repository.CreateAssetHistory(&model.AssetHistory{
		ID: uuid.New().String(), ChangeType: model.Depreciated,
		ChangedAt:   time.Now(),
		Description: fmt.Sprintf("Monthly depreciation %.2f (%s). New value: %d", monthly, method, newValue),
		AssetID:     assetID, UserID: userID,
	})

	return &dto.DepreciationResult{
		AssetID: asset.ID, AssetName: asset.AssetName,
		AcquisitionPrice: asset.AcquisitionPrice, CurrentValue: newValue,
		MonthlyAmount: monthly, Method: method,
		UsefulLifeMonths: asset.UsefulLifeMonths,
	}, nil
}

func RevalueAsset(assetID string, newValue int, reason, userID string) error {
	asset, err := repository.FindAssetByID(assetID)
	if err != nil {
		return fmt.Errorf("asset not found")
	}
	oldVal := asset.CurrentValue
	asset.CurrentValue = newValue
	if err := repository.UpdateAsset(asset); err != nil {
		return err
	}
	repository.CreateAssetHistory(&model.AssetHistory{
		ID: uuid.New().String(), ChangeType: model.Revalued,
		ChangedAt:   time.Now(),
		Description: fmt.Sprintf("Revalued from %d to %d. Reason: %s", oldVal, newValue, reason),
		AssetID:     assetID, UserID: userID,
	})
	return nil
}
