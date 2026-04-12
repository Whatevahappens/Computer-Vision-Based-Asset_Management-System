package model

type AssetModel struct {
	ID                      string
	Brand                   string
	ModelName               string
	AssetModelType          AssetType
	Category                AssetCategory
	DefaultUnitPrice        int
	DefaultUsefulLifeMonths int
	DepreciationMethod      DepreciationMethod
}
