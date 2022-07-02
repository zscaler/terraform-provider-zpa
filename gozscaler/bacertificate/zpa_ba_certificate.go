package bacertificate

import (
	"fmt"
	"net/http"

	"github.com/zscaler/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfigV1                = "/mgmtconfig/v1/admin/customers/"
	baCertificateEndpoint       = "/clientlessCertificate"
	mgmtConfigV2                = "/mgmtconfig/v2/admin/customers/"
	baCertificateIssuedEndpoint = "/clientlessCertificate/issued"
)

type BaCertificate struct {
	CName               string   `json:"cName,omitempty"`
	CertChain           string   `json:"certChain,omitempty"`
	CreationTime        string   `json:"creationTime,omitempty"`
	Description         string   `json:"description,omitempty"`
	ID                  string   `json:"id,omitempty"`
	IssuedBy            string   `json:"issuedBy,omitempty"`
	IssuedTo            string   `json:"issuedTo,omitempty"`
	ModifiedBy          string   `json:"modifiedBy,omitempty"`
	ModifiedTime        string   `json:"modifiedTime,omitempty"`
	Name                string   `json:"name,omitempty"`
	San                 []string `json:"san,omitempty"`
	SerialNo            string   `json:"serialNo,omitempty"`
	Status              string   `json:"status,omitempty"`
	ValidFromInEpochSec string   `json:"validFromInEpochSec,omitempty"`
	ValidToInEpochSec   string   `json:"validToInEpochSec,omitempty"`
}

type GenerateCSR struct {
	CSRString   string   `json:"csrString,omitempty"`
	Description string   `json:"description,omitempty"`
	Name        string   `json:"name,omitempty"`
	SANS        []string `json:"sans,omitempty"`
	Subject     string   `json:"subject,omitempty"`
}

func (service *Service) Get(baCertificateID string) (*BaCertificate, *http.Response, error) {
	v := new(BaCertificate)
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfigV1+service.Client.Config.CustomerID+baCertificateEndpoint, baCertificateID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetIssuedByName(CertName string) (*BaCertificate, *http.Response, error) {
	var v struct {
		List []BaCertificate `json:"list"`
	}
	relativeURL := fmt.Sprintf(mgmtConfigV2 + service.Client.Config.CustomerID + baCertificateIssuedEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: CertName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, baCertificate := range v.List {
		if baCertificate.Name == CertName {
			return &baCertificate, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no issued certificate named '%s' was found", CertName)
}

func (service *Service) Create(certificate BaCertificate) (*BaCertificate, *http.Response, error) {
	v := new(BaCertificate)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfigV1+service.Client.Config.CustomerID+baCertificateEndpoint, nil, certificate, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// /zpn/api/v1/admin/customers/{customerId}/clientlessCertificate/generateCSR
// Generate a certificate request
func (service *Service) Reorder(policySetID, ruleId string, order int) (*http.Response, error) {
	path := fmt.Sprintf("/zpn/api/v1/admin/customers/%s/generateCSR", service.Client.Config.CustomerID)
	resp, err := service.Client.NewPrivateRequestDo("PUT", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// /zpn/api/v1/admin/customers/{customerId}/clientlessCertificate/{certificateId}
// Update name/description on ClientlessCertificate
func (service *Service) Update(id string, certificate BaCertificate) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.Config.CustomerID+baCertificateEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, certificate, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.Config.CustomerID+baCertificateEndpoint, id)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
