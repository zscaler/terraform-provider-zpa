package policysetcontroller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig = "/mgmtconfig/v1/admin/customers/"
)

type PolicySet struct {
	CreationTime string       `json:"creationTime,omitempty"`
	Description  string       `json:"description,omitempty"`
	Enabled      bool         `json:"enabled"`
	ID           string       `json:"id,omitempty"`
	ModifiedBy   string       `json:"modifiedBy,omitempty"`
	ModifiedTime string       `json:"modifiedTime,omitempty"`
	Name         string       `json:"name,omitempty"`
	PolicyType   string       `json:"policyType,omitempty"`
	Rules        []PolicyRule `json:"rules"`
}
type PolicyRule struct {
	Action                   string               `json:"action,omitempty"`
	ActionID                 string               `json:"actionId,omitempty"`
	BypassDefaultRule        bool                 `json:"bypassDefaultRule"`
	CreationTime             string               `json:"creationTime,omitempty"`
	CustomMsg                string               `json:"customMsg,omitempty"`
	DefaultRule              bool                 `json:"defaultRule,omitempty"`
	DefaultRuleName          string               `json:"defaultRuleName,omitempty"`
	Description              string               `json:"description,omitempty"`
	ID                       string               `json:"id,omitempty"`
	IsolationDefaultRule     bool                 `json:"isolationDefaultRule"`
	ModifiedBy               string               `json:"modifiedBy,omitempty"`
	ModifiedTime             string               `json:"modifiedTime,omitempty"`
	Name                     string               `json:"name,omitempty"`
	Operator                 string               `json:"operator,omitempty"`
	PolicySetID              string               `json:"policySetId"`
	PolicyType               string               `json:"policyType,omitempty"`
	Priority                 string               `json:"priority,omitempty"`
	ReauthDefaultRule        bool                 `json:"reauthDefaultRule"`
	ReauthIdleTimeout        string               `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout            string               `json:"reauthTimeout,omitempty"`
	RuleOrder                string               `json:"ruleOrder"`
	LSSDefaultRule           bool                 `json:"lssDefaultRule"`
	ZpnCbiProfileID          string               `json:"zpnCbiProfileId,omitempty"`
	ZpnInspectionProfileID   string               `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName string               `json:"zpnInspectionProfileName,omitempty"`
	Conditions               []Conditions         `json:"conditions,omitempty"`
	AppServerGroups          []AppServerGroups    `json:"appServerGroups,omitempty"`
	AppConnectorGroups       []AppConnectorGroups `json:"appConnectorGroups,omitempty"`
}

type Conditions struct {
	CreationTime string     `json:"creationTime,omitempty"`
	ID           string     `json:"id,omitempty"`
	ModifiedBy   string     `json:"modifiedBy,omitempty"`
	ModifiedTime string     `json:"modifiedTime,omitempty"`
	Negated      bool       `json:"negated"`
	Operands     []Operands `json:"operands,omitempty"`
	Operator     string     `json:"operator,omitempty"`
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
	ID string `json:"id,omitempty"`
}
type AppConnectorGroups struct {
	ID string `json:"id,omitempty"`
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

// GET --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule/{ruleId}
func (service *Service) GetPolicyRule(policySetID, ruleId string) (*PolicyRule, *http.Response, error) {
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

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/policySet/{policySetId}/rule/{ruleId}/reorder/{newOrder}
func (service *Service) Reorder(policySetID, ruleId string, order int) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/rule/%s/reorder/%d", policySetID, ruleId, order)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) RulesCount() (int, *http.Response, error) {
	v := new(Count)
	relativeURL := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/rules/policyType/GLOBAL_POLICY/count", service.Client.Config.CustomerID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return 0, nil, err
	}
	count, err := strconv.Atoi(v.Count)
	return count, resp, err
}
<<<<<<< HEAD
=======

// Get the global policy. This API will be deprecated in a future release.
// GET /mgmtconfig/v1/admin/customers/{customerId}/policySet/global
func (service *Service) GetPolicySet() (*PolicySet, *http.Response, error) {
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

/*
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
*/
>>>>>>> master
