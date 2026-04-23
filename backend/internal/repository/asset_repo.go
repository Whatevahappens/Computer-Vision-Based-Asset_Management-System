package repository

import (
	"backend/internal/database"
	"backend/internal/model"
)

func CreateAsset(asset *model.Asset) error {
	return database.DB.Create(asset).Error
}

func FindAssetByID(id string) (*model.Asset, error) {
	var asset model.Asset
	err := database.DB.Preload("AssetModel").Preload("Department").
		Preload("Location").Preload("AssingnedUser").
		First(&asset, "id = ?", id).Error
	return &asset, err
}

func FindAssetByBarcode(barcode string) (*model.Asset, error) {
	var asset model.Asset
	err := database.DB.First(&asset, "barcode = ?", barcode).Error
	return &asset, err
}

func ListAssets(offset, limit int, search, status string) ([]model.Asset, int64, error) {
	var assets []model.Asset
	var total int64
	q := database.DB.Model(&model.Asset{})
	if search != "" {
		q = q.Where("asset_name ILIKE ? OR barcode ILIKE ? OR serial_number ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Count(&total)
	err := q.Preload("AssetModel").Preload("Location").Preload("Department").Preload("AssignedUser").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&assets).Error
	return assets, total, err
}

func ListAssetsByUser(userID string) ([]model.Asset, error) {
	var assets []model.Asset
	err := database.DB.Preload("Location").Where("assigned_user_id = ?", userID).Find(&assets).Error
	return assets, err
}

func ListAssetsByLocation(locationID string) ([]model.Asset, error) {
	var assets []model.Asset
	err := database.DB.Where("location_id = ? AND status = ?", locationID, model.AssetActive).Find(&assets).Error
	return assets, err
}

func CountAssetsByLocationAndCategory()
