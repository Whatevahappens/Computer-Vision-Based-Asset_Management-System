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

func CountAssetsByLocationAndCategory(locationID string) (map[model.AssetCategory]int, error) {
	type result struct {
		Category model.AssetCategory
		Count    int
	}
	var results []result
	err := database.DB.Model(&model.Asset{}).
		Select("COALESCE((SELECT category FROM asset_models WHERE id = asset_model_id), 'OTHER') as category, COUNT(*) as count").
		Where("location_id = ? AND status = ?", locationID, model.AssetActive).
		Group("category").Find(&results).Error
	if err != nil {
		return nil, err
	}
	m := make(map[model.AssetCategory]int)
	for _, r := range results {
		m[r.Category] = r.Count
	}
	return m, nil
}

func UpdateAsset(asset *model.Asset) error {
	return database.DB.Save(asset).Error
}

func CountAssets() (int64, error) {
	var count int64
	err := database.DB.Model(&model.Asset{}).Count(&count).Error
	return count, err
}

func CountActiveAssets() (int64, error) {
	var count int64
	err := database.DB.Model(&model.Asset{}).Where("status = ?", model.AssetActive).Count(&count).Error
	return count, err
}

func SumAssetValues() (int64, error) {
	var total int64
	err := database.DB.Model(&model.Asset{}).Select("COALESCE(SUM(current_value), 0)").Scan(&total).Error
	return total, err
}

func GetAllAssetsForReport() ([]model.Asset, error) {
	var assets []model.Asset
	err := database.DB.Preload("AssetModel").Preload("Location").Preload("Department").Preload("AssignedUser").
		Order("created_at DESC").Find(&assets).Error
	return assets, err
}

func CreateAssetHistory(h *model.AssetHistory) error {
	return database.DB.Create(h).Error
}

func GetAssetHistory(assetID string) ([]model.AssetHistory, error) {
	var history []model.AssetHistory
	err := database.DB.Where("asset_id = ?", assetID).Order("changed_at DESC").Find(&history).Error
	return history, err
}

func CreateAssetModel(m *model.AssetModel) error {
	return database.DB.Create(m).Error
}

func ListAssetModels() ([]model.AssetModel, error) {
	var models []model.AssetModel
	err := database.DB.Order("brand, model_name").Find(&models).Error
	return models, err
}

func FindAssetModelByID(id string) (*model.AssetModel, error) {
	var m model.AssetModel
	err := database.DB.First(&m, "id = ?", id).Error
	return &m, err
}

func CreateLocation(l *model.Location) error {
	return database.DB.Create(l).Error
}

func ListLocations() ([]model.Location, error) {
	var locs []model.Location
	err := database.DB.Order("name").Find(&locs).Error
	return locs, err
}

func FindLocationByID(id string) (*model.Location, error) {
	var l model.Location
	err := database.DB.First(&l, "id = ?", id).Error
	return &l, err
}

func UpdateLocation(l *model.Location) error {
	return database.DB.Save(l).Error
}

func DeleteLocation(id string) error {
	return database.DB.Delete(&model.Location{}, "id = ?", id).Error
}

func CreateDepartment(d *model.Department) error {
	return database.DB.Create(d).Error
}

func ListDepartments() ([]model.Department, error) {
	var deps []model.Department
	err := database.DB.Order("name").Find(&deps).Error
	return deps, err
}

func FindDepartmentByID(id string) (*model.Department, error) {
	var d model.Department
	err := database.DB.First(&d, "id = ?", id).Error
	return &d, err
}

func UpdateDepartment(d *model.Department) error {
	return database.DB.Save(d).Error
}

func DeleteDepartment(id string) error {
	return database.DB.Delete(&model.Department{}, "id = ?", id).Error
}

func CreateDocument(doc *model.Document) error {
	return database.DB.Create(doc).Error
}

func ListDocumentByAsset(assetID string) ([]model.Document, error) {
	var docs []model.Document
	err := database.DB.Where("asset_id = ?", assetID).Find(&docs).Error
	return docs, err
}
