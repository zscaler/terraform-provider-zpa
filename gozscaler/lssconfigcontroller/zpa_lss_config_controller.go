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
	Config          *Config            `json:"config"`
	ID              string             `json:"id"`
	ConnectorGroups *[]ConnectorGroups `json:"connectorGroups"`
	PolicyRule      *PolicyRule        `json:"policyRule"`
	// PolicyRuleResource PolicyRuleResource `json:"policyRuleResource"`
}
type Config struct {
	AuditMessage  string   `json:"auditMessage"`
	CreationTime  string   `json:"creationTime"`
	Description   string   `json:"description"`
	Enabled       bool     `json:"enabled"`
	Filter        []string `json:"filter"`
	Format        string   `json:"format"`
	ID            string   `json:"id"`
	ModifiedBy    string   `json:"modifiedBy"`
	ModifiedTime  string   `json:"modifiedTime"`
	Name          string   `json:"name"`
	LssHost       string   `json:"lssHost"`
	LssPort       string   `json:"lssPort"`
	SourceLogType string   `json:"sourceLogType"`
	UseTLS        bool     `json:"useTls"`
}
type ConnectorGroups struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PolicyRule struct {
	Action                   string       `json:"action"`
	ActionID                 string       `json:"actionId"`
	BypassDefaultRule        bool         `json:"bypassDefaultRule"`
	CreationTime             string       `json:"creationTime"`
	CustomMsg                string       `json:"customMsg"`
	DefaultRule              bool         `json:"defaultRule"`
	Description              string       `json:"description"`
	ID                       string       `json:"id"`
	IsolationDefaultRule     bool         `json:"isolationDefaultRule"`
	ModifiedBy               string       `json:"modifiedBy"`
	ModifiedTime             string       `json:"modifiedTime"`
	Name                     string       `json:"name"`
	Operator                 string       `json:"operator"`
	PolicySetID              string       `json:"policySetId"`
	PolicyType               string       `json:"policyType"`
	Priority                 string       `json:"priority"`
	ReauthDefaultRule        bool         `json:"reauthDefaultRule"`
	ReauthIdleTimeout        string       `json:"reauthIdleTimeout"`
	ReauthTimeout            string       `json:"reauthTimeout"`
	RuleOrder                string       `json:"ruleOrder"`
	LssDefaultRule           bool         `json:"lssDefaultRule"`
	ZpnCbiProfileID          string       `json:"zpnCbiProfileId"`
	ZpnInspectionProfileID   string       `json:"zpnInspectionProfileId"`
	ZpnInspectionProfileName string       `json:"zpnInspectionProfileName"`
	Conditions               []Conditions `json:"conditions"`
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
