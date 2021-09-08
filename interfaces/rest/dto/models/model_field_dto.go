package models

type ModelFieldAddParamsApiDto struct {
	Type          string `json:"type" url:"type"`
	ReferenceCode string `json:"referenceCode" url:"referenceCode"`
	ModelCode     string `json:"modelCode" url:"modelCode"`
	Multiple      bool   `json:"multiple" url:"multiple"`
}
