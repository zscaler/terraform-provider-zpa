package common

const (
	DefaultPageSize = 500
)

type Pagination struct {
	PageSize int    `json:"pagesize" url:"pagesize"`
	Page     int    `json:"page,omitempty" url:"page"`
	Search   string //`json:"search,omitempty" url:"search"`
}

type NetworkPorts struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}
