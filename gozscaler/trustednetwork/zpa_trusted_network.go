package trustednetwork

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig             = "/mgmtconfig/v2/admin/customers/"
	trustedNetworkEndpoint = "/network"
)

type TrustedNetwork struct {
	CreationTime     string `json:"creationTime,omitempty"`
	Domain           string `json:"domain,omitempty"`
	ID               string `json:"id,omitempty"`
	MasterCustomerID string `json:"masterCustomerId,omitempty"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedTime     string `json:"modifiedTime,omitempty"`
	Name             string `json:"name,omitempty"`
	NetworkID        string `json:"networkId,omitempty"`
	ZscalerCloud     string `json:"zscalerCloud,omitempty"`
}

func (service *Service) Get(networkID string) (*TrustedNetwork, *http.Response, error) {
	v := new(TrustedNetwork)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+trustedNetworkEndpoint, networkID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByNetID(netID string) (*TrustedNetwork, *http.Response, error) {
	var v struct {
		List []TrustedNetwork `json:"list"`
	}
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + trustedNetworkEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, trustedNetwork := range v.List {
		if trustedNetwork.NetworkID == netID {
			return &trustedNetwork, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no trusted network with NetworkID '%s' was found", netID)
}

func (service *Service) GetByName(trustedNetworkName string) (*TrustedNetwork, *http.Response, error) {
	var v struct {
		List []TrustedNetwork `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + trustedNetworkEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: trustedNetworkName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, trustedNetwork := range v.List {
		if strings.EqualFold(trustedNetwork.Name, trustedNetworkName) {
			return &trustedNetwork, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no trusted network named '%s' was found", trustedNetworkName)
}
