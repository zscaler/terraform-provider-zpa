package inspection_profile

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig                = "/mgmtconfig/v1/admin/customers/"
	inspectionProfileEndpoint = "/inspectionProfile"
)

type InspectionProfile struct {
	ID                                string                    `json:"id,omitempty"`
	CommonGlobalOverrideActionsConfig map[string]interface{}    `json:"commonGlobalOverrideActionsConfig,omitempty"`
	CreationTime                      string                    `json:"creationTime,omitempty"`
	Description                       string                    `json:"description,omitempty"`
	GlobalControlActions              []string                  `json:"globalControlActions,omitempty"`
	IncarnationNumber                 string                    `json:"incarnationNumber,omitempty"`
	ModifiedBy                        string                    `json:"modifiedBy,omitempty"`
	ModifiedTime                      string                    `json:"modifiedTime,omitempty"`
	Name                              string                    `json:"name,omitempty"`
	ParanoiaLevel                     string                    `json:"paranoiaLevel,omitempty"`
	PredefinedControlsVersion         string                    `json:"predefinedControlsVersion,omitempty"`
	ControlInfoResource               []ControlInfoResource     `json:"controlsInfo,omitempty"`
	CustomControls                    []InspectionCustomControl `json:"customControls,omitempty"`
	PredefinedControls                []PredefinedControls      `json:"predefinedControls,omitempty"`
}

type ControlInfoResource struct {
	ControlType string `json:"controlType,omitempty"`
	Count       string `json:"count,omitempty"`
}

type InspectionCustomControl struct {
	Action                           string                   `json:"action,omitempty"`
	ActionValue                      string                   `json:"actionValue,omitempty"`
	AssociatedInspectionProfileNames []AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
	Rules                            []Rules                  `json:"rules,omitempty"`
	ControlNumber                    string                   `json:"controlNumber,omitempty"`
	ControlRuleJson                  string                   `json:"controlRuleJson,omitempty"`
	CreationTime                     string                   `json:"creationTime,omitempty"`
	DefaultAction                    string                   `json:"defaultAction,omitempty"`
	DefaultActionValue               string                   `json:"defaultActionValue,omitempty"`
	Description                      string                   `json:"description,omitempty"`
	ID                               string                   `json:"id,omitempty"`
	ModifiedBy                       string                   `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                   `json:"modifiedTime,omitempty"`
	Name                             string                   `json:"name,omitempty"`
	ParanoiaLevel                    string                   `json:"paranoiaLevel,omitempty"`
	Severity                         string                   `json:"severity,omitempty"`
	Type                             string                   `json:"type,omitempty"`
	Version                          string                   `json:"version,omitempty"`
}

type Rules struct {
	Conditions []Conditions `json:"conditions,omitempty"`
	Names      string       `json:"names,omitempty"`
	Type       string       `json:"type,omitempty"`
}

type Conditions struct {
	LHS string `json:"lhs,omitempty"`
	OP  string `json:"op,omitempty"`
	RHS string `json:"rhs,omitempty"`
}

type PredefinedControls struct {
	ID                               string                   `json:"id,omitempty"`
	Name                             string                   `json:"name,omitempty"`
	Action                           string                   `json:"action,omitempty"`
	ActionValue                      string                   `json:"actionValue,omitempty"`
	AssociatedInspectionProfileNames []AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
	Attachment                       string                   `json:"attachment,omitempty"`
	ControlGroup                     string                   `json:"controlGroup,omitempty"`
	ControlNumber                    string                   `json:"controlNumber,omitempty"`
	CreationTime                     string                   `json:"creationTime,omitempty"`
	DefaultAction                    string                   `json:"defaultAction,omitempty"`
	DefaultActionValue               string                   `json:"defaultActionValue,omitempty"`
	Description                      string                   `json:"description,omitempty"`
	ModifiedBy                       string                   `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                   `json:"modifiedTime,omitempty"`
	ParanoiaLevel                    string                   `json:"paranoiaLevel,omitempty"`
	Severity                         string                   `json:"severity,omitempty"`
	Version                          string                   `json:"version,omitempty"`
}

type AssociatedProfileNames struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (service *Service) Get(profileID string) (*InspectionProfile, *http.Response, error) {
	v := new(InspectionProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+inspectionProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(profileName string) (*InspectionProfile, *http.Response, error) {
	var v struct {
		List []InspectionProfile `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + inspectionProfileEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: profileName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, inspection := range v.List {
		if strings.EqualFold(inspection.Name, profileName) {
			return &inspection, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no inspection profile named '%s' was found", profileName)
}

func (service *Service) Create(inspectionProfile InspectionProfile) (*InspectionProfile, *http.Response, error) {
	v := new(InspectionProfile)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+inspectionProfileEndpoint, nil, inspectionProfile, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Update(profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+inspectionProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, nil, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *Service) PutAssociate(profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+inspectionProfileEndpoint, profileID+"associateAllPredefinedControls")
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, nil, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *Service) PutDeassociate(profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+inspectionProfileEndpoint, profileID+"associateAllPredefinedControls")
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, nil, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *Service) Patch(profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+inspectionProfileEndpoint, profileID+"patch")
	resp, err := service.Client.NewRequestDo("PATCH", relativeURL, nil, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *Service) Delete(profileID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+inspectionProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
