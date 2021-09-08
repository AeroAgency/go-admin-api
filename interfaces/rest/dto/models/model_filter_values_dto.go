package models

type ModelFilterValuesParamsDto struct {
	ModelCode               string `json:"model_code"`
	ModelFieldCode          string `json:"model_field_code"`
	ModelFieldType          string `json:"model_field_type"`
	ModelFieldModelCode     string `json:"model_field_model_code"`
	ModelFieldReferenceCode string `json:"model_field_reference_code"`
	Limit                   int    `json:"limit"`
	Offset                  int    `json:"offset"`
	Query                   string `json:"query"`
	Multiple                bool   `json:"multiple" url:"multiple"`
}

type ValueRows struct {
	Items []ValueRow `json:"items"`
	Total int        `json:"total"`
}

type ValueRow struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}
