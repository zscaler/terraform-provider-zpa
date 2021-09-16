package trustednetwork

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	mgmtConfig             = "/mgmtconfig/v1/admin/customers/"
	trustedNetworkEndpoint = "/network"
)

type TrustedNetwork struct {
	CreationTime string `json:"creationTime,omitempty"`
	Domain       string `json:"domain,omitempty"`
	ID           string `json:"id,omitempty"`
	ModifiedBy   string `json:"modifiedBy,omitempty"`
	ModifiedTime string `json:"modifiedTime,omitempty"`
	Name         string `json:"name,omitempty"`
	NetworkID    string `json:"networkId,omitempty"`
	ZscalerCloud string `json:"zscalerCloud,omitempty"`
}

func (service *Service) Get(networkId string) (*TrustedNetwork, *http.Response, error) {
	v := new(TrustedNetwork)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+trustedNetworkEndpoint, networkId)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(name string) (*TrustedNetwork, *http.Response, error) {
	var v []TrustedNetwork
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + trustedNetworkEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, struct{ pagesize int }{
		pagesize: 500,
	}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, trustedNetwork := range v {
		if strings.EqualFold(trustedNetwork.Name, name) {
			return &trustedNetwork, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no saml trusted network named '%s' was found", name)
}
