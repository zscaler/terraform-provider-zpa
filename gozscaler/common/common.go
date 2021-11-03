package common

type Pagination struct {
	PageSize int    `json:"pagesize"`
	Page     int    `json:"page,omitempty"`
	Search   string `json:"search,omitempty"`
}
