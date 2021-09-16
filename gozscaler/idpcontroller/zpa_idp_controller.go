package idpcontroller

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	mgmtConfig                 = "/mgmtconfig/v1/admin/customers/"
	idpControllerGroupEndpoint = "/idp"
)

type IdpController struct {
	AutoProvision               string         `json:"autoProvision,omitempty"`
	CreationTime                string         `json:"creationTime,omitempty"`
	Description                 string         `json:"description,omitempty"`
	DisableSamlBasedPolicy      bool           `json:"disableSamlBasedPolicy,omitempty"`
	Domainlist                  []string       `json:"domainList,omitempty"`
	EnableScimBasedPolicy       bool           `json:"enableScimBasedPolicy,omitempty"`
	Enabled                     bool           `json:"enabled,omitempty"`
	ID                          string         `json:"id,omitempty"`
	IdpEntityID                 string         `json:"idpEntityId,omitempty"`
	LoginNameAttribute          string         `json:"loginNameAttribute,omitempty"`
	LoginURL                    string         `json:"loginUrl,omitempty"`
	ModifiedBy                  string         `json:"modifiedBy,omitempty"`
	ModifiedTime                string         `json:"modifiedTime,omitempty"`
	Name                        string         `json:"name,omitempty"`
	ReauthOnUserUpdate          bool           `json:"reauthOnUserUpdate,omitempty"`
	RedirectBinding             bool           `json:"redirectBinding,omitempty"`
	ScimEnabled                 bool           `json:"scimEnabled,omitempty"`
	ScimServiceProviderEndpoint string         `json:"scimServiceProviderEndpoint,omitempty"`
	ScimSharedSecret            string         `json:"scimSharedSecret,omitempty"`
	ScimSharedSecretExists      bool           `json:"scimSharedSecretExists,omitempty"`
	SignSamlRequest             string         `json:"signSamlRequest,,omitempty"`
	SsoType                     []string       `json:"ssoType,omitempty"`
	UseCustomSpMetadata         bool           `json:"useCustomSPMetadata,omitempty"`
	AdminMetadata               AdminMetadata  `json:"adminMetadata,omitempty"`
	UserMetadata                UserMetadata   `json:"userMetadata,omitempty"`
	Certificates                []Certificates `json:"certificates"`
}

type AdminMetadata struct {
	CertificateURL string `json:"certificateUrl"`
	SpEntityID     string `json:"spEntityId"`
	SpMetadataURL  string `json:"spMetadataUrl"`
	SpPostURL      string `json:"spPostUrl"`
}
type Certificates struct {
	Cname          string `json:"cName,omitempty"`
	Certificate    string `json:"certificate,omitempty"`
	SerialNo       string `json:"serialNo,omitempty"`
	ValidFrominSec string `json:"validFromInSec,omitempty"`
	ValidToinSec   string `json:"validToInSec,omitempty"`
}
type UserMetadata struct {
	CertificateURL string `json:"certificateUrl,omitempty"`
	SpEntityID     string `json:"spEntityId,omitempty"`
	SpMetadataURL  string `json:"spMetadataUrl,omitempty"`
	SpPostURL      string `json:"spPostUrl,omitempty"`
}

func (service *Service) Get(IdpID string) (*IdpController, *http.Response, error) {
	v := new(IdpController)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+idpControllerGroupEndpoint, IdpID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(name string) (*IdpController, *http.Response, error) {
	var v []IdpController
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + idpControllerGroupEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, struct{ pagesize int }{
		pagesize: 500,
	}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, idpController := range v {
		if strings.EqualFold(idpController.Name, name) {
			return &idpController, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no Idp-Controller named '%s' was found", name)
}
