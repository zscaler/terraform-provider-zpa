package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredentialpool"
)

func dataSourcePRACredentialPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePRACredentialPoolRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the privileged credential pool",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the privileged credential pool",
			},
			"credential_mapping_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"credential_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol type that was designated for that particular privileged credential. The protocol type options are SSH, RDP, and VNC. Each protocol type has its own credential requirements.",
			},
			"credentials": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "List of PRA Credentials",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
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

func dataSourcePRACredentialPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *pracredentialpool.CredentialPool
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for pra credential pool %s\n", id)
		res, _, err := pracredentialpool.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for  pra credential poolr name %s\n", name)
		res, _, err := pracredentialpool.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("credential_mapping_count", resp.CredentialMappingCount)
		_ = d.Set("credential_type", resp.CredentialType)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)
		_ = d.Set("credentials", flattenCredentials(resp.PRACredentials))
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any pra credential pool with name '%s' or id '%s'", name, id))
	}

	return nil
}
