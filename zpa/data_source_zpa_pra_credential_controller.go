package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/pracredential"
)

func dataSourcePRACredentialController() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePRACredentialControllerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the privileged credential",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the privileged credential",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the privileged credential",
			},
			"last_credential_reset_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the privileged credential was last reset",
			},
			"credential_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol type that was designated for that particular privileged credential. The protocol type options are SSH, RDP, and VNC. Each protocol type has its own credential requirements.",
			},
			"user_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain name associated with the username. You can also include the domain name as part of the username. The domain name only needs to be specified with logging in to an RDP console that is connected to an Active Directory Domain.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username for the login you want to use for the privileged credential",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The password associated with the username for the login you want to use for the privileged credential",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the privileged credential is created",
			},
			"modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the tenant who modified the privileged credential",
			},
			"modified_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the privileged credential is modified",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.",
			},
			"microtenant_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Microtenant",
			},
		},
	}
}

func dataSourcePRACredentialControllerRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PRACredential

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}
	var resp *pracredential.Credential
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for credential controller %s\n", id)
		res, _, err := pracredential.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for credential controller name %s\n", name)
		res, _, err := pracredential.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("last_credential_reset_time", resp.LastCredentialResetTime)
		_ = d.Set("credential_type", resp.CredentialType)
		_ = d.Set("user_domain", resp.UserDomain)
		_ = d.Set("username", resp.UserName)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)
	} else {
		return fmt.Errorf("couldn't find any credential controller with name '%s' or id '%s'", name, id)
	}

	return nil
}
