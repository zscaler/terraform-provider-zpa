package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/enrollmentcert"
)

func dataSourceEnrollmentCert() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEnrollmentCertRead,
		Schema: map[string]*schema.Schema{
			"allow_signing": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "certificate text in pem format.",
			},
			"client_cert_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"csr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"issued_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issued_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifiedby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_cert_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_cert_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key_present": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"serial_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_from_in_epoch_sec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_to_in_epoch_sec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zrsa_encrypted_private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zrsa_encrypted_session_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEnrollmentCertRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *enrollmentcert.EnrollmentCert
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for signing certificate %s\n", id)
		res, _, err := zClient.enrollmentcert.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for signing certificate name %s\n", name)
		res, _, err := zClient.enrollmentcert.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("allow_signing", resp.AllowSigning)
		_ = d.Set("cname", resp.Cname)
		_ = d.Set("certificate", resp.Certificate)
		_ = d.Set("client_cert_type", resp.ClientCertType)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("csr", resp.CSR)
		_ = d.Set("description", resp.Description)
		_ = d.Set("issued_by", resp.IssuedBy)
		_ = d.Set("issued_to", resp.IssuedTo)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("parent_cert_id", resp.ParentCertID)
		_ = d.Set("parent_cert_name", resp.ParentCertName)
		_ = d.Set("private_key", resp.PrivateKey)
		_ = d.Set("private_key_present", resp.PrivateKeyPresent)
		_ = d.Set("serial_no", resp.SerialNo)
		_ = d.Set("valid_from_in_epoch_sec", resp.ValidFromInEpochSec)
		_ = d.Set("valid_to_in_epoch_sec", resp.ValidToInEpochSec)
		_ = d.Set("zrsa_encrypted_private_key", resp.ZrsaEncryptedPrivateKey)
		_ = d.Set("zrsa_encrypted_session_key", resp.ZrsaEncryptedSessionKey)
	} else {
		return fmt.Errorf("couldn't find any signing certificate with name '%s' or id '%s'", name, id)
	}

	return nil
}
