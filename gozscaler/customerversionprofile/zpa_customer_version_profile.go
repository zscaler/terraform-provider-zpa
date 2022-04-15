package customerversionprofile

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig                     = "/mgmtconfig/v1/admin/customers/"
	customerVersionProfileEndpoint = "/visible/versionProfiles"
)

type CustomerVersionProfile struct {
	CreationTime                  string                        `json:"creationTime"`
	CustomScopeCustomerIDs        []CustomScopeCustomerIDs      `json:"customScopeCustomerIds"`
	CustomScopeRequestCustomerIDs CustomScopeRequestCustomerIDs `json:"customScopeRequestCustomerIds"`
	CustomerID                    string                        `json:"customerId"`
	Description                   string                        `json:"description"`
	ID                            string                        `json:"id"`
	ModifiedBy                    string                        `json:"modifiedBy"`
	ModifiedTime                  string                        `json:"modifiedTime"`
	Name                          string                        `json:"name"`
	UpgradePriority               string                        `json:"upgradePriority"`
	Versions                      []Versions                    `json:"versions"`
	VisibilityScope               string                        `json:"visibilityScope"`
}

type CustomScopeCustomerIDs struct {
	CustomerID           string `json:"customerId"`
	ExcludeConstellation bool   `json:"excludeConstellation"`
	Name                 string `json:"name"`
}
type CustomScopeRequestCustomerIDs struct {
	AddCustomerIDs    string `json:"addCustomerIds"`
	DeletecustomerIDs string `json:"deleteCustomerIds"`
}
type Versions struct {
	CreationTime             string `json:"creationTime"`
	CustomerID               string `json:"customerId"`
	ID                       string `json:"id"`
	ModifiedBy               string `json:"modifiedBy"`
	ModifiedTime             string `json:"modifiedTime"`
	Platform                 string `json:"platform"`
	RestartAfterUptimeInDays string `json:"restartAfterUptimeInDays"`
	Role                     string `json:"role"`
	Version                  string `json:"version"`
	VersionProfileGID        string `json:"version_profile_gid"`
}

func (service *Service) Get(versionID string) (*CustomerVersionProfile, *http.Response, error) {
	v := new(CustomerVersionProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+customerVersionProfileEndpoint, versionID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(versionProfileName string) (*CustomerVersionProfile, *http.Response, error) {
	var v struct {
		List []CustomerVersionProfile `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + customerVersionProfileEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: versionProfileName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range v.List {
		if strings.EqualFold(app.Name, versionProfileName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no version profile named '%s' was found", versionProfileName)
}
