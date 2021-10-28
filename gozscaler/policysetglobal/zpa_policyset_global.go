package policysetglobal

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	mgmtConfig = "/mgmtconfig/v1/admin/customers/"
)

type PolicySet struct {
	CreationTime string  `json:"creationTime,omitempty"`
	Description  string  `json:"description,omitempty"`
	Enabled      bool    `json:"enabled"`
	ID           string  `json:"id,omitempty"`
	ModifiedBy   string  `json:"modifiedBy,omitempty"`
	ModifiedTime string  `json:"modifiedTime,omitempty"`
	Name         string  `json:"name,omitempty"`
	PolicyType   string  `json:"policyType,omitempty"`
	Rules        []Rules `json:"rules"`
}

type Rules struct {
	Action                   string       `json:"action,omitempty"`
	ActionID                 string       `json:"actionId,omitempty"`
	BypassDefaultRule        bool         `json:"bypassDefaultRule"`
	CreationTime             string       `json:"creationTime,omitempty"`
	CustomMsg                string       `json:"customMsg,omitempty"`
	Description              string       `json:"description,omitempty"`
	ID                       string       `json:"id,omitempty"`
	IsolationDefaultRule     bool         `json:"isolationDefaultRule"`
	ModifiedBy               string       `json:"modifiedBy,omitempty"`
	ModifiedTime             string       `json:"modifiedTime,omitempty"`
	Name                     string       `json:"name,omitempty"`
	Operator                 string       `json:"operator,omitempty"`
	PolicySetID              string       `json:"policySetId,omitempty"`
	PolicyType               string       `json:"policyType,omitempty"`
	Priority                 string       `json:"priority,omitempty"`
	ReauthDefaultRule        bool         `json:"reauthDefaultRule"`
	ReauthIdleTimeout        string       `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout            string       `json:"reauthTimeout,omitempty"`
	RuleOrder                string       `json:"ruleOrder,omitempty"`
	ZpnCbiProfileID          string       `json:"zpnCbiProfileId,omitempty"`
	ZpnInspectionProfileID   string       `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName string       `json:"zpnInspectionProfileName,omitempty"`
	Conditions               []Conditions `json:"conditions,omitempty"`
}
type Conditions struct {
	CreationTime string      `json:"creationTime,omitempty"`
	ID           string      `json:"id,omitempty"`
	ModifiedBy   string      `json:"modifiedBy,omitempty"`
	ModifiedTime string      `json:"modifiedTime,omitempty"`
	Negated      bool        `json:"negated"`
	Operands     *[]Operands `json:"operands,omitempty"`
	Operator     string      `json:"operator,omitempty"`
}
type Operands struct {
	CreationTime string `json:"creationTime,omitempty"`
	ID           string `json:"id,omitempty"`
	IdpID        string `json:"idpId,omitempty"`
	LHS          string `json:"lhs,omitempty"`
	ModifiedBy   string `json:"modifiedBy,omitempty"`
	ModifiedTime string `json:"modifiedTime,omitempty"`
	Name         string `json:"name,omitempty"`
	ObjectType   string `json:"objectType,omitempty"`
	RHS          string `json:"rhs,omitempty"`
}
type AppServerGroups struct {
	ConfigSpace      string `json:"configSpace,omitempty"`
	CreationTime     string `json:"creationTime,omitempty"`
	Description      string `json:"description,omitempty"`
	Enabled          bool   `json:"enabled"`
	ID               string `json:"id,omitempty"`
	DynamicDiscovery bool   `json:"dynamicDiscovery"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedTime     string `json:"modifiedTime,omitempty"`
	Name             string `json:"name,omitempty"`
}

type Count struct {
	Count string `json:"count"`
}

func (service *Service) GetByPolicyType(policyType string) (*PolicySet, *http.Response, error) {
	v := new(PolicySet)
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + "/policySet/policyType/" + policyType)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Get the global policy. This API will be deprecated in a future release.
// GET /mgmtconfig/v1/admin/customers/{customerId}/policySet/global
func (service *Service) Get() (*PolicySet, *http.Response, error) {
	v := new(PolicySet)
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + "/policySet/global")
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Get the authentication policy and all rules for a Timeout policy rule. This API will be deprecated in a future release.
// /mgmtconfig/v1/admin/customers/{customerId}/policySet/reauth
func (service *Service) GetReauth() (*PolicySet, *http.Response, error) {
	v := new(PolicySet)
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + "/policySet/reauth")
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Get the bypass policy and all rules for a Client Forwarding policy rule. This API will be deprecated in a future release.
// GET mgmtconfig/v1/admin/customers/{customerId}/policySet/bypass
func (service *Service) GetBypass() (*PolicySet, *http.Response, error) {
	v := new(PolicySet)
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + "/policySet/bypass")
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) RulesCount() (int, *http.Response, error) {
	v := new(Count)
	relativeURL := fmt.Sprintf("/zpn/api/v1/admin/customers/%s/policySet/rules/policyType/GLOBAL_POLICY/count", service.Client.Config.CustomerID)
	resp, err := service.Client.NewPrivateRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return 0, nil, err
	}
	count, err := strconv.Atoi(v.Count)
	return count, resp, err
}
