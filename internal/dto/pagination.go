package dto

type Pagination struct {
	CurrentPage  int64 `json:"current_page"`
	TotalPages   int64 `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	HasMore      bool  `json:"has_more"`
	Limit        int64 `json:"limit"`
}
