package models

type EmptyObj struct{}

type PagingResponse struct {
	TotalPage    int         `json:"total_page"`
	TotalRecords int64       `json:"total_records"`
	CurrentPage  int         `json:"current_page"`
	Data         interface{} `json:"data"`
}
