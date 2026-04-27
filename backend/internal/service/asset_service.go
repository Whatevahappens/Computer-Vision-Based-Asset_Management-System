package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateBarcode() string {
	now := time.Now()
	return fmt.Sprintf("BC-%d-%04d", now.Year(), now.UnixNano()%10000)
}

func CreateAsset(name, serial string, price int, acqDate time.Time, lifeMonths int, nature, desc, modelID, deptID, locID, userID string) (*model.Asset, error) {
	asset := &model.Asset{
		ID:               uuid.New().String(),
		Barcode:          GenerateBarcode(),
		SerialNumber:     serial,
		AssetName:        name,
		AcquisitionPrice: price,
		AcquisitionDate:  acqDate,
		UsefulLifeMonths: lifeMonths,
		CurrentValue:     price,
		Status:           model.AssetActive,
		Nature:           model.AssetNature(nature),
		Description:      desc,
	}
	if modelID != "" {
		asset.AssetModelID = &modelID
	}
	if deptID != "" {
		asset.DepartmentID = &deptID
	}
	if locID != "" {
		asset.LocationID = &locID
	}
	if err := repository.CreateAsset(asset); err != nil {
		return nil, err
	}

	repository.CreateAssetHistory(&model.AssetHistory{
		ID: uuid.New().String(), ChangeType: model.Created,
		ChangedAt: time.Now(), Description: "Asset created",
		AssetID: asset.ID, UserID: userID,
	})

	return asset, nil
}

func AssignAsset(assetID, toUserID, locationID, notes, performedBy string) error {
	asset, err := repository.FindAssetByID(assetID)
	if err != nil {
		return fmt.Errorf("asset not found")
	}
	asset.AssignedUserID = &toUserID
	if locationID != "" {
		asset.LocationID = &locationID
	}
	if err := repository.UpdateAsset(asset); err != nil {
		return err
	}
	repository.CreateAssetHistory(&model.AssetHistory{
		ID: uuid.New().String(), ChangeType: model.Assigned,
		ChangedAt: time.Now(), Description: "Assigned to user " + toUserID + ". " + notes,
		AssetID: assetID, UserID: performedBy,
	})
	repository.CreateNotification(&model.Notification{
		ID: uuid.New().String(), Title: "Эд хөрөнгө хуваарилагдлаа",
		Message: asset.AssetName + " танд хуваарилагдлаа",
		Type:    model.Info, UserID: toUserID, CreatedAt: time.Now(),
	})
	return nil
}

func TransferAsset(assetID, toUserID, locationID, notes, performedBy string) error {
	asset, err := repository.FindAssetByID(assetID)
	if err != nil {
		return fmt.Errorf("asset not found")
	}
	oldUserID := ""
	if asset.AssignedUserID != nil {
		oldUserID = *asset.AssignedUserID
	}
	asset.AssignedUserID = &toUserID
	if locationID != "" {
		asset.LocationID = &locationID
	}
	if err := repository.UpdateAsset(asset); err != nil {
		return err
	}
	repository.CreateAssetHistory(&model.AssetHistory{
		ID: uuid.New().String(), ChangeType: model.Transferred,
		ChangedAt: time.Now(), Description: fmt.Sprintf("Transferred from %s to %s. %s", oldUserID, toUserID, notes),
		AssetID: assetID, UserID: performedBy,
	})

	if oldUserID != "" {
		repository.CreateNotification(&model.Notification{
			ID: uuid.New().String(), Title: "Эд хөрөнгө буцаагдлаа",
			Message: asset.AssetName + " буцаан авагдлаа.",
			Type:    model.Info, UserID: oldUserID, CreatedAt: time.Now(),
		})
	}
	repository.CreateNotification(&model.Notification{
		ID: uuid.New().String(), Title: "Эд хөрөнгө хуваарилагдлаа",
		Message: asset.AssetName + " танд шилжүүлэгдлээ.",
		Type:    model.Info, UserID: toUserID, CreatedAt: time.Now(),
	})
	return nil
}

func DisposeAsset(assetID, reason string, residualValue int, notes, performedBy string) error {
	asset, err := repository.FindAssetByID(assetID)
	if err != nil {
		return fmt.Errorf("asset not found")
	}
	asset.Status = model.AssetDisposed
	asset.CurrentValue = residualValue
	asset.AssignedUserID = nil
	if err := repository.UpdateAsset(asset); err != nil {
		return err
	}
	repository.CreateAssetHistory(&model.AssetHistory{
		ID: uuid.New().String(), ChangeType: model.Disposed,
		ChangedAt: time.Now(), Description: fmt.Sprintf("Disposed. Reason: %s, %s", reason, notes),
		AssetID: assetID, UserID: performedBy,
	})
	return nil
}
