package cloudconnectorgroup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig                  = "/mgmtconfig/v1/admin/customers/"
	cloudConnectorGroupEndpoint = "/cloudConnectorGroup"
)

type CloudConnectorGroup struct {
	CreationTime    string            `json:"creationTime,omitempty"`
	Description     string            `json:"description,omitempty"`
	CloudConnectors []CloudConnectors `json:"cloudConnectors,omitempty"`
	Enabled         bool              `json:"enabled,omitempty"`
	GeolocationID   string            `json:"geoLocationId,omitempty"`
	ID              string            `json:"id,omitempty"`
	ModifiedBy      string            `json:"modifiedBy,omitempty"`
	ModifiedTime    string            `json:"modifiedTime,omitempty"`
	Name            string            `json:"name,omitempty"`
	ZiaCloud        string            `json:"ziaCloud,omitempty"`
	ZiaOrgid        string            `json:"ziaOrgId,omitempty"`
}
type CloudConnectors struct {
	CreationTime string                 `json:"creationTime,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Enabled      bool                   `json:"enabled,omitempty"`
	Fingerprint  string                 `json:"fingerprint,omitempty"`
	ID           string                 `json:"id,omitempty"`
	IPACL        []string               `json:"ipAcl,omitempty"`
	IssuedCertID string                 `json:"issuedCertId,omitempty"`
	ModifiedBy   string                 `json:"modifiedBy,omitempty"`
	ModifiedTime string                 `json:"modifiedTime,omitempty"`
	SigningCert  map[string]interface{} `json:"signingCert,omitempty"`
	Name         string                 `json:"name,omitempty"`
}

func (service *Service) Get(cloudConnectorGroupID string) (*CloudConnectorGroup, *http.Response, error) {
	v := new(CloudConnectorGroup)
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + cloudConnectorGroupEndpoint + "/" + cloudConnectorGroupID
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(cloudConnectorGroupName string) (*CloudConnectorGroup, *http.Response, error) {
	var v struct {
		List []CloudConnectorGroup `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + cloudConnectorGroupEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: cloudConnectorGroupName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range v.List {
		if strings.EqualFold(app.Name, cloudConnectorGroupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application named '%s' was found", cloudConnectorGroupName)
}
