package entity

type Pagination struct {
	Page           int64 `json:"page"`
	PageSize       int64 `json:"page_size"`
	TotalDocuments int64 `json:"total_documents"`
}
