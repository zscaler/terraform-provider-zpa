package inspection_predefined_controls

import (
	"fmt"
	"log"
	"net/http"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig           = "/mgmtconfig/v1/admin/customers/"
	predControlsEndpoint = "/inspectionControls/predefined"
)

type PredefinedControls struct {
	ID                               string                          `json:"id,omitempty"`
	Name                             string                          `json:"name,omitempty"`
	Action                           string                          `json:"action,omitempty"`
	ActionValue                      string                          `json:"actionValue,omitempty"`
	AssociatedInspectionProfileNames []common.AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
	Attachment                       string                          `json:"attachment,omitempty"`
	ControlGroup                     string                          `json:"controlGroup,omitempty"`
	ControlNumber                    string                          `json:"controlNumber,omitempty"`
	CreationTime                     string                          `json:"creationTime,omitempty"`
	DefaultAction                    string                          `json:"defaultAction,omitempty"`
	DefaultActionValue               string                          `json:"defaultActionValue,omitempty"`
	Description                      string                          `json:"description,omitempty"`
	ModifiedBy                       string                          `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                          `json:"modifiedTime,omitempty"`
	ParanoiaLevel                    string                          `json:"paranoiaLevel,omitempty"`
	Severity                         string                          `json:"severity,omitempty"`
	Version                          string                          `json:"version,omitempty"`
}

/*
type AssociatedProfileNames struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
*/

// Get Predefined Controls by ID
// https://help.zscaler.com/zpa/api-reference#/inspection-control-controller/getPredefinedControlById
func (service *Service) Get(controlID string) (*PredefinedControls, *http.Response, error) {
	v := new(PredefinedControls)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+predControlsEndpoint, controlID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetAll() (*PredefinedControls, *http.Response, error) {
	v := new(PredefinedControls)
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + predControlsEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("[INFO] got predefined controls:%#v", v)
	return v, resp, nil
}

/*
// Get All Predefined Inspection Controls - Need to create multiple search criteria
//https://help.zscaler.com/zpa/api-reference#/inspection-control-controller/getAllInspectionControls
func (service *Service) GetAll() (*PredefinedControls, *http.Response, error) {
	var v struct {
		List []PredefinedControls `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + predControlsEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: controlName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, inspection := range v.List {
		if strings.EqualFold(inspection.Name) {
			return &inspection, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no predefined control named '%s' was found", controlName)
}
*/
