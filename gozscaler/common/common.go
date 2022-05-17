package common

const (
	DefaultPageSize = 500
)

type Pagination struct {
	PageSize int    `json:"pagesize" url:"pagesize"`
	Page     int    `json:"page,omitempty" url:"page,omitempty"`
	Search   string `json:"-" url:"-"`
	Search2  string `json:"search,omitempty" url:"search,omitempty"`
}

type NetworkPorts struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// ZPA Inspection Rules
type Rules struct {
	Conditions []Conditions `json:"conditions,omitempty"`
	Names      string       `json:"names,omitempty"`
	Type       string       `json:"type,omitempty"`
	Version    string       `json:"version,omitempty"`
}

type Conditions struct {
	LHS string `json:"lhs,omitempty"`
	OP  string `json:"op,omitempty"`
	RHS string `json:"rhs,omitempty"`
}
type AssociatedProfileNames struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
