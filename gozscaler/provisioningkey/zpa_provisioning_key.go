package provisioningkey

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig = "/mgmtconfig/v1/admin/customers/"
)

type ProvisioningKey struct {
	AppConnectorGroupID   string   `json:"appConnectorGroupId,omitempty"`
	AppConnectorGroupName string   `json:"appConnectorGroupName,omitempty"`
	CreationTime          string   `json:"creationTime,omitempty"`
	Enabled               bool     `json:"enabled,omitempty"`
	ExpirationInEpochSec  string   `json:"expirationInEpochSec,omitempty"`
	ID                    string   `json:"id,omitempty"`
	IPACL                 []string `json:"ipAcl,omitempty"`
	MaxUsage              string   `json:"maxUsage,omitempty"`
	ModifiedBy            string   `json:"modifiedBy,omitempty"`
	ModifiedTime          string   `json:"modifiedTime,omitempty"`
	Name                  string   `json:"name,omitempty"`
	ProvisioningKey       string   `json:"provisioningKey,omitempty"`
	EnrollmentCertID      string   `json:"enrollmentCertId,omitempty"`
	EnrollmentCertName    string   `json:"enrollmentCertName,omitempty"`
	UIConfig              string   `json:"uiConfig,omitempty"`
	UsageCount            string   `json:"usageCount,omitempty"`
	ZcomponentID          string   `json:"zcomponentId,omitempty"`
	ZcomponentName        string   `json:"zcomponentName,omitempty"`
}

// GET --> mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey
func (service *Service) Get(associationType, provisioningKeyID string) (*ProvisioningKey, *http.Response, error) {
	v := new(ProvisioningKey)
	url := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/associationType/%s/provisioningKey/%s", associationType, provisioningKeyID)
	resp, err := service.Client.NewRequestDo("GET", url, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// GET --> mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey
func (service *Service) GetByName(associationType, name string) (*ProvisioningKey, *http.Response, error) {
	list, resp, err := service.GetAll(associationType)
	if err != nil {
		return nil, nil, err
	}
	for _, provisioningKey := range list {
		if strings.EqualFold(provisioningKey.Name, name) {
			return &provisioningKey, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no ProvisioningKey named '%s' was found", name)
}

// GET --> mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey
func (service *Service) GetAll(associationType string) ([]ProvisioningKey, *http.Response, error) {
	var v struct {
		List []ProvisioningKey `json:"list"`
	}
	url := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/associationType/%s/provisioningKey", associationType)
	resp, err := service.Client.NewRequestDo("GET", url, common.Pagination{PageSize: 500}, nil, &v)
	return v.List, resp, err
}

// POST --> /mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey
func (service *Service) Create(associationType string, provisioningKey *ProvisioningKey) (*ProvisioningKey, *http.Response, error) {
	v := new(ProvisioningKey)
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/associationType/%s/provisioningKey", associationType)
	resp, err := service.Client.NewRequestDo("POST", path, nil, provisioningKey, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey/{provisioningKeyId}
func (service *Service) Update(associationType, provisioningKeyID string, provisioningKey *ProvisioningKey) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/associationType/%s/provisioningKey/%s", associationType, provisioningKeyID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, provisioningKey, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// DELETE --> /mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey/{provisioningKeyId}
func (service *Service) Delete(associationType, provisioningKeyID string) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/associationType/%s/provisioningKey/%s", associationType, provisioningKeyID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
