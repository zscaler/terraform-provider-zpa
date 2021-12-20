package common

const (
	DefaultPageSize = 500
)

type Pagination struct {
	PageSize int    `json:"pagesize" url:"pagesize"`
	Page     int    `json:"page,omitempty" url:"page"`
<<<<<<< HEAD
	Search   string //`json:"search,omitempty" url:"search"`
=======
	Search   string `json:"-" url:"-"`
>>>>>>> 615856b1a9082873acb1096923c7c71d4d91a39f
}
