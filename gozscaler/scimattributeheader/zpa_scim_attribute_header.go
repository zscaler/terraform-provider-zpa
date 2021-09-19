package scimattributeheader

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	mgmtConfig       = "/mgmtconfig/v1/admin/customers/"
	idpId            = "/idp"
	scimAttrEndpoint = "/scimattribute"
)

type ScimAttributeHeader struct {
	CanonicalValues []string `json:"canonicalValues,omitempty"`
	CaseSensitive   bool     `json:"caseSensitive,omitempty"`
	CreationTime    string   `json:"creationTime,omitempty,"`
	DataType        string   `json:"dataType,omitempty"`
	Description     string   `json:"description,omitempty"`
	ID              string   `json:"id,omitempty"`
	IdpID           string   `json:"idpId,omitempty"`
	ModifiedBy      string   `json:"modifiedBy,omitempty"`
	ModifiedTime    string   `json:"modifiedTime,omitempty"`
	MultiValued     bool     `json:"multivalued,omitempty"`
	Mutability      string   `json:"mutability,omitempty"`
	Name            string   `json:"name,omitempty"`
	Required        bool     `json:"required,omitempty"`
	Returned        string   `json:"returned,omitempty"`
	SchemaURI       string   `json:"schemaURI,omitempty"`
	Uniqueness      bool     `json:"uniqueness,omitempty"`
}

func (service *Service) Get(ScimAttrHeaderID string) (*ScimAttributeHeader, *http.Response, error) {
	v := new(ScimAttributeHeader)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+idpId+scimAttrEndpoint, ScimAttrHeaderID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(scimAttributeName, IdpId string) (*ScimAttributeHeader, *http.Response, error) {
	var v struct {
		List []ScimAttributeHeader `json:"list"`
	}
	relativeURL := fmt.Sprintf("%s/%s%s", mgmtConfig+service.Client.Config.CustomerID+idpId, IdpId, scimAttrEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, struct{ pagesize int }{
		pagesize: 500,
	}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, scimAttribute := range v.List {
		if strings.EqualFold(scimAttribute.Name, scimAttributeName) {
			return &scimAttribute, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no scim named '%s' was found", scimAttributeName)
}
