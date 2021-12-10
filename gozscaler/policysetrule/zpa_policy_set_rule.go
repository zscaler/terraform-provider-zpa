package policysetrule

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig = "/mgmtconfig/v1/admin/customers/"
)

type PolicyRule struct {
	Action             string               `json:"action,omitempty"`
	ActionID           string               `json:"actionId,omitempty"`
	BypassDefaultRule  bool                 `json:"bypassDefaultRule"`
	CustomMsg          string               `json:"customMsg,omitempty"`
	DefaultRule        bool                 `json:"defaultRule,omitempty"`
	Description        string               `json:"description,omitempty"`
	ID                 string               `json:"id,omitempty"`
	Name               string               `json:"name,omitempty"`
	Operator           string               `json:"operator,omitempty"`
	PolicySetID        string               `json:"policySetId"`
	PolicyType         string               `json:"policyType,omitempty"`
	Priority           string               `json:"priority,omitempty"`
	ReauthDefaultRule  bool                 `json:"reauthDefaultRule"`
	ReauthIdleTimeout  string               `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout      string               `json:"reauthTimeout,omitempty"`
	RuleOrder          string               `json:"ruleOrder"`
	LSSDefaultRule     bool                 `json:"lssDefaultRule"`
	Conditions         []Conditions         `json:"conditions,omitempty"`
	AppServerGroups    []AppServerGroups    `json:"appServerGroups,omitempty"`
	AppConnectorGroups []AppConnectorGroups `json:"appConnectorGroups,omitempty"`
}

type Conditions struct {
	ID       string     `json:"id,omitempty"`
	Negated  bool       `json:"negated"`
	Operands []Operands `json:"operands,omitempty"`
	Operator string     `json:"operator,omitempty"`
}
type Operands struct {
	ID         string `json:"id,omitempty"`
	IdpID      string `json:"idpId,omitempty"`
	LHS        string `json:"lhs,omitempty"`
	ObjectType string `json:"objectType,omitempty"`
	RHS        string `json:"rhs,omitempty"`
	Name       string `json:"name,omitempty"`
}

type AppServerGroups struct {
	ID string `json:"id,omitempty"`
}
type AppConnectorGroups struct {
	ID string `json:"id,omitempty"`
}

// GET --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func (service *Service) Get(policySetID, ruleId string) (*PolicyRule, *http.Response, error) {
	v := new(PolicyRule)
	url := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("GET", url, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// POST --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule
func (service *Service) Create(rule *PolicyRule) (*PolicyRule, *http.Response, error) {
	v := new(PolicyRule)
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/rule", rule.PolicySetID)
	resp, err := service.Client.NewRequestDo("POST", path, nil, &rule, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// PUT --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func (service *Service) Update(policySetID, ruleId string, policySetRule *PolicyRule) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, policySetRule, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/policySet/{policySetId}/rule/{ruleId}/reorder/{newOrder}
func (service *Service) Reorder(policySetID, ruleId string, order int) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/rule/%s/reorder/%d", policySetID, ruleId, order)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// DELETE --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func (service *Service) Delete(policySetID, ruleId string) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetByNameAndType(policyType, ruleName string) (*PolicyRule, *http.Response, error) {
	var v struct {
		List []PolicyRule `json:"list"`
	}
	url := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/rules/policyType/%s", policyType)
	resp, err := service.Client.NewRequestDo("GET", url, common.Pagination{PageSize: common.DefaultPageSize, Search: ruleName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, p := range v.List {
		if strings.EqualFold(ruleName, p.Name) {
			return &p, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no policy rule named :%s found", ruleName)
}

func (service *Service) GetByNameAndTypes(policyTypes []string, ruleName string) (p *PolicyRule, resp *http.Response, err error) {
	for _, policyType := range policyTypes {
		p, resp, err = service.GetByNameAndType(policyType, ruleName)
		if err != nil {
			continue
		} else {
			return
		}
	}
	return
}
