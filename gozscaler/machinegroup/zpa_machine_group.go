package machinegroup

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	mgmtConfig                  = "/mgmtconfig/v1/admin/customers/"
	machineGroupEndpoint string = "/machineGroup"
)

type MachineGroup struct {
	CreationTime string     `json:"creationTime,omitempty"`
	Description  string     `json:"description,omitempty"`
	Enabled      bool       `json:"enabled,omitempty"`
	ID           string     `json:"id,omitempty"`
	Machines     []Machines `json:"machines,omitempty"`
	ModifiedBy   string     `json:"modifiedBy,omitempty"`
	ModifiedTime string     `json:"modifiedTime,omitempty"`
	Name         string     `json:"name,omitempty"`
}

type Machines struct {
	CreationTime     string                 `json:"creationTime,omitempty"`
	Description      string                 `json:"description,omitempty"`
	Fingerprint      string                 `json:"fingerprint,omitempty"`
	ID               string                 `json:"id,omitempty"`
	IssuedCertID     string                 `json:"issuedCertId,omitempty"`
	MachineGroupID   string                 `json:"machineGroupId,omitempty"`
	MachineGroupName string                 `json:"machineGroupName,omitempty"`
	MachineTokenID   string                 `json:"machineTokenId,omitempty"`
	ModifiedBy       string                 `json:"modifiedBy,omitempty"`
	ModifiedTime     string                 `json:"modifiedTime,omitempty"`
	Name             string                 `json:"name,omitempty"`
	SigningCert      map[string]interface{} `json:"signingCert,omitempty"`
}

func (service *Service) Get(machineGroupID string) (*MachineGroup, *http.Response, error) {
	v := new(MachineGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+machineGroupEndpoint, machineGroupID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(machineGroupName string) (*MachineGroup, *http.Response, error) {
	var v struct {
		List []MachineGroup `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + machineGroupEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, struct{ pagesize int }{
		pagesize: 500,
	}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range v.List {
		if strings.EqualFold(app.Name, machineGroupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application named '%s' was found", machineGroupName)
}
