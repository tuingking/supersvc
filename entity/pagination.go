package entity

type Pagination struct {
	Page    int64 `json:"page"`
	Size    int64 `json:"size"`
	Total   int64 `json:"total"`
	HasNext bool  `json:"has_next"`
}
