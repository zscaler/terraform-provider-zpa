package lssconfigcontroller

import (
	"fmt"
	"net/http"
)

const (
	mgmtConfig        = "/mgmtconfig/v2/admin/customers/"
	lssConfigEndpoint = "/lssConfig"
)

type LSSResource struct {
	LSSConfig          *LSSConfig          `json:"config"`
	ID                 string              `json:"id"`
	ConnectorGroups    []ConnectorGroups   `json:"connectorGroups,omitempty"`
	PolicyRule         *PolicyRule         `json:"policyRule,omitempty"`
	PolicyRuleResource *PolicyRuleResource `json:"policyRuleResource,omitempty"`
}
type LSSConfig struct {
	AuditMessage  string   `json:"auditMessage,omitempty"`
	CreationTime  string   `json:"creationTime,omitempty"`
	Description   string   `json:"description,omitempty"`
	Enabled       bool     `json:"enabled,omitempty"`
	Filter        []string `json:"filter,omitempty"`
	Format        string   `json:"format,omitempty"`
	ID            string   `json:"id,omitempty"`
	ModifiedBy    string   `json:"modifiedBy,omitempty"`
	ModifiedTime  string   `json:"modifiedTime,omitempty"`
	Name          string   `json:"name,omitempty"`
	LSSHost       string   `json:"lssHost,omitempty"`
	LSSPort       string   `json:"lssPort,omitempty"`
	SourceLogType string   `json:"sourceLogType,omitempty"`
	UseTLS        bool     `json:"useTls,omitempty"`
}
type ConnectorGroups struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type PolicyRuleResource struct {
	Action                   string                         `json:"action,omitempty"`
	ActionID                 string                         `json:"actionId,omitempty"`
	BypassDefaultRule        bool                           `json:"bypassDefaultRule,omitempty"`
	CreationTime             string                         `json:"creationTime,omitempty"`
	CustomMsg                string                         `json:"customMsg,omitempty"`
	DefaultRule              bool                           `json:"defaultRule,omitempty"`
	Description              string                         `json:"description,omitempty"`
	ID                       string                         `json:"id,omitempty"`
	IsolationDefaultRule     bool                           `json:"isolationDefaultRule,omitempty"`
	ModifiedBy               string                         `json:"modifiedBy,omitempty"`
	ModifiedTime             string                         `json:"modifiedTime,omitempty"`
	Name                     string                         `json:"name,omitempty"`
	Operator                 string                         `json:"operator,omitempty"`
	PolicySetID              string                         `json:"policySetId,omitempty"`
	PolicyType               string                         `json:"policyType,omitempty"`
	Priority                 string                         `json:"priority,omitempty"`
	ReauthDefaultRule        bool                           `json:"reauthDefaultRule,omitempty"`
	ReauthIdleTimeout        string                         `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout            string                         `json:"reauthTimeout,omitempty"`
	RuleOrder                string                         `json:"ruleOrder,omitempty"`
	LssDefaultRule           bool                           `json:"lssDefaultRule,omitempty"`
	ZpnCbiProfileID          string                         `json:"zpnCbiProfileId,omitempty"`
	ZpnInspectionProfileID   string                         `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName string                         `json:"zpnInspectionProfileName,omitempty"`
	Conditions               []PolicyRuleResourceConditions `json:"conditions,omitempty"`
}
type PolicyRule struct {
	Action                   string       `json:"action,omitempty"`
	ActionID                 string       `json:"actionId,omitempty"`
	BypassDefaultRule        bool         `json:"bypassDefaultRule,omitempty"`
	CreationTime             string       `json:"creationTime,omitempty"`
	CustomMsg                string       `json:"customMsg,omitempty"`
	DefaultRule              bool         `json:"defaultRule,omitempty"`
	Description              string       `json:"description,omitempty"`
	ID                       string       `json:"id,omitempty"`
	IsolationDefaultRule     bool         `json:"isolationDefaultRule,omitempty"`
	ModifiedBy               string       `json:"modifiedBy,omitempty"`
	ModifiedTime             string       `json:"modifiedTime,omitempty"`
	Name                     string       `json:"name,omitempty"`
	Operator                 string       `json:"operator,omitempty"`
	PolicySetID              string       `json:"policySetId,omitempty"`
	PolicyType               string       `json:"policyType,omitempty"`
	Priority                 string       `json:"priority,omitempty"`
	ReauthDefaultRule        bool         `json:"reauthDefaultRule,omitempty"`
	ReauthIdleTimeout        string       `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout            string       `json:"reauthTimeout,omitempty"`
	RuleOrder                string       `json:"ruleOrder,omitempty"`
	LssDefaultRule           bool         `json:"lssDefaultRule,omitempty"`
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

type PolicyRuleResourceConditions struct {
	CreationTime string                        `json:"creationTime,omitempty"`
	ID           string                        `json:"id,omitempty"`
	ModifiedBy   string                        `json:"modifiedBy,omitempty"`
	ModifiedTime string                        `json:"modifiedTime,omitempty"`
	Negated      bool                          `json:"negated"`
	Operands     *[]PolicyRuleResourceOperands `json:"operands,omitempty"`
	Operator     string                        `json:"operator,omitempty"`
}
type PolicyRuleResourceOperands struct {
	ObjectType string   `json:"objectType,omitempty"`
	Values     []string `json:"values,omitempty"`
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

func (service *Service) Get(lssID string) (*LSSResource, *http.Response, error) {
	v := new(LSSResource)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+lssConfigEndpoint, lssID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Create(lssConfig *LSSResource) (*LSSResource, *http.Response, error) {
	v := new(LSSResource)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+lssConfigEndpoint, nil, lssConfig, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(lssID string, lssConfig *LSSResource) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+lssConfigEndpoint, lssID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, lssConfig, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(lssID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+lssConfigEndpoint, lssID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
