package scimuserattributeheader

import (
	"fmt"
	"net/http"
)

const (
	userConfig           = "/userconfig/v1/customers/"
	idpId                = "/idp"
	scimUserAttrEndpoint = "/scimattribute"
)

type ScimUserAttributeHeader struct {
	List []string `json:"list,omitempty"`
}

func (service *Service) GetAll() (*ScimUserAttributeHeader, *http.Response, error) {
	v := new(ScimUserAttributeHeader)
	relativeURL := fmt.Sprintf("%s", userConfig+service.Client.Config.CustomerID+idpId+scimUserAttrEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// func (service *Service) GetAttributeByName(scimAttributeName, IdpId string) (*ScimUserAttributeHeader, *http.Response, error) {
// 	var v struct {
// 		List []ScimUserAttributeHeader `json:"list"`
// 	}
// 	relativeURL := fmt.Sprintf("%s/%s%s", userConfig+service.Client.Config.CustomerID+idpId, IdpId, scimUserAttrEndpoint)
// 	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: scimAttributeName}, nil, &v)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	for _, scimAttribute := range v.List {
// 		if strings.EqualFold(scimAttribute.Name, scimAttributeName) {
// 			return &scimAttribute, resp, nil
// 		}
// 	}
// 	return nil, resp, fmt.Errorf("no scim named '%s' was found", scimAttributeName)
// }
