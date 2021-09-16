package models

// SENDING TO EXTERNAL API
type ModelElementsListParamsApiDto struct {
	ModelCode string                                `json:"modelCode"`
	Offset    int                                   `json:"offset"`
	Limit     int                                   `json:"limit"`
	Sort      string                                `json:"sort"`
	Order     string                                `json:"order"`
	Select    []string                              `json:"select"`
	Filter    []ModelElementsListParamsApiDtoFilter `json:"filter"`
}

type ModelElementsListParamsApiDtoFilter FieldValue

type FieldValue struct {
	Code   string   `json:"code"`
	Values []string `json:"values"`
}

// INNER
type ModelElements struct {
	Items []ModelElement `json:"items"`
	Total int            `json:"total"`
}

type ModelElement struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ModelElementDetail map[string]string

type ModelElementLink struct {
	LinkId string `json:"link_id"`
}

type ModelElementCreateApiDto struct {
	Id          string
	FieldValues []FieldValue `json:"fieldValues"`
}

type ModelElementIdDto struct {
	ModelElementId string `json:"modelElementId"`
}
