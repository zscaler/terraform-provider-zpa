package idpcontroller

import (
	"fmt"
	"net/http"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig            = "/mgmtconfig/v2/admin/customers/"
	mgmtConfigV1          = "/mgmtconfig/v1/admin/customers/"
	idpControllerEndpoint = "/idp"
)

type IdpController struct {
	AdminSpSigningCertID        string        `json:"adminSpSigningCertId,omitempty"`
	AutoProvision               string        `json:"autoProvision,omitempty"`
	CreationTime                string        `json:"creationTime,omitempty"`
	Description                 string        `json:"description,omitempty"`
	DisableSamlBasedPolicy      bool          `json:"disableSamlBasedPolicy"`
	Domainlist                  []string      `json:"domainList,omitempty"`
	EnableScimBasedPolicy       bool          `json:"enableScimBasedPolicy"`
	Enabled                     bool          `json:"enabled"`
	ID                          string        `json:"id,omitempty"`
	IdpEntityID                 string        `json:"idpEntityId,omitempty"`
	LoginNameAttribute          string        `json:"loginNameAttribute,omitempty"`
	LoginURL                    string        `json:"loginUrl,omitempty"`
	ModifiedBy                  string        `json:"modifiedBy,omitempty"`
	ModifiedTime                string        `json:"modifiedTime,omitempty"`
	Name                        string        `json:"name,omitempty"`
	ReauthOnUserUpdate          bool          `json:"reauthOnUserUpdate"`
	RedirectBinding             bool          `json:"redirectBinding"`
	ScimEnabled                 bool          `json:"scimEnabled"`
	ScimServiceProviderEndpoint string        `json:"scimServiceProviderEndpoint,omitempty"`
	ScimSharedSecret            string        `json:"scimSharedSecret,omitempty"`
	ScimSharedSecretExists      bool          `json:"scimSharedSecretExists,omitempty"`
	SignSamlRequest             string        `json:"signSamlRequest,,omitempty"`
	SsoType                     []string      `json:"ssoType,omitempty"`
	UseCustomSpMetadata         bool          `json:"useCustomSPMetadata"`
	UserSpSigningCertID         string        `json:"userSpSigningCertId,omitempty"`
	AdminMetadata               AdminMetadata `json:"adminMetadata,omitempty"`
	UserMetadata                UserMetadata  `json:"userMetadata,omitempty"`
}

type AdminMetadata struct {
	CertificateURL string `json:"certificateUrl"`
	SpBaseURL      string `json:"spBaseUrl"`
	SpEntityID     string `json:"spEntityId"`
	SpMetadataURL  string `json:"spMetadataUrl"`
	SpPostURL      string `json:"spPostUrl"`
}
type UserMetadata struct {
	CertificateURL string `json:"certificateUrl,omitempty"`
	SpBaseURL      string `json:"spBaseUrl"`
	SpEntityID     string `json:"spEntityId,omitempty"`
	SpMetadataURL  string `json:"spMetadataUrl,omitempty"`
	SpPostURL      string `json:"spPostUrl,omitempty"`
}

func (service *Service) Get(IdpID string) (*IdpController, *http.Response, error) {
	v := new(IdpController)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.Config.CustomerID+idpControllerEndpoint, IdpID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(idpName string) (*IdpController, *http.Response, error) {
	var v struct {
		List []IdpController `json:"list"`
	}
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + idpControllerEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{
		PageSize: 500,
		Search:   idpName,
	}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, idpController := range v.List {
		if idpController.Name == idpName {
			return &idpController, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no Idp-Controller named '%s' was found", idpName)
}
