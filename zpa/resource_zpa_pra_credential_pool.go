package zpa

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredentialpool"
)

func resourcePRACredentialPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePRACredentialPoolCreate,
		ReadContext:   resourcePRACredentialPoolRead,
		UpdateContext: resourcePRACredentialPoolUpdate,
		DeleteContext: resourcePRACredentialPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				microTenantID := GetString(d.Get("microtenant_id"))
				if microTenantID != "" {
					service = service.WithMicroTenant(microTenantID)
				}

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := pracredentialpool.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(resp.ID)
						_ = d.Set("id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the privileged credential",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the privileged credential",
			},
			"credential_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The protocol type that was designated for that particular privileged credential. The protocol type options are SSH, RDP, and VNC. Each protocol type has its own credential requirements.",
				ValidateFunc: validation.StringInSlice([]string{
					"USERNAME_PASSWORD",
					"SSH_KEY",
					"PASSWORD",
				}, false),
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
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant.",
			},
		},
	}
}

func resourcePRACredentialPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandPRACredentialPool(d)

	sanitizeFields(&req)
	log.Printf("[INFO] Creating pra credential pool with request\n%+v\n", req)

	credController, _, err := pracredentialpool.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created pra credential pool request. ID: %v\n", credController)

	d.SetId(credController.ID)
	return resourcePRACredentialPoolRead(ctx, d, meta)
}

func resourcePRACredentialPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := pracredentialpool.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing pra credential pool %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting pra credential pool:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("credential_type", resp.CredentialType)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("credentials", flattenCredentials(resp.PRACredentials))

	return nil
}

func resourcePRACredentialPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("credential_type") {
		oldType, newType := d.GetChange("credential_type")
		return diag.FromErr(fmt.Errorf("changing 'credential_type' from '%s' to '%s' is not allowed", oldType, newType))
	}

	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating pra credential pool ID: %v\n", id)

	req := expandPRACredentialPool(d)

	if _, _, err := pracredentialpool.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := pracredentialpool.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourcePRACredentialPoolRead(ctx, d, meta)
}

func resourcePRACredentialPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting pra credential pool ID: %v\n", d.Id())

	// Detach the pra credential from all policy rules before attempting to delete it
	if err := detachPRACredentialFromPolicy(ctx, d.Id(), service); err != nil {
		return diag.FromErr(fmt.Errorf("error detaching pra credential with ID %s from PolicySetControllers: %s", d.Id(), err))
	}

	if _, err := pracredentialpool.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] pra credential pool deleted")
	return nil
}

func expandPRACredentialPool(d *schema.ResourceData) pracredentialpool.CredentialPool {
	credPool := pracredentialpool.CredentialPool{
		ID:             d.Id(),
		Name:           d.Get("name").(string),
		CredentialType: d.Get("credential_type").(string),
		MicroTenantID:  d.Get("microtenant_id").(string),
		PRACredentials: expandCredentials(d),
	}
	return credPool
}

func expandCredentials(d *schema.ResourceData) []common.CommonIDName {
	credentialsSet, ok := d.GetOk("credentials")
	if !ok {
		return nil
	}

	credentials := []common.CommonIDName{}
	for _, credRaw := range credentialsSet.(*schema.Set).List() {
		credMap := credRaw.(map[string]interface{})
		if idSet, ok := credMap["id"].(*schema.Set); ok && idSet.Len() > 0 {
			for _, id := range idSet.List() {
				credentials = append(credentials, common.CommonIDName{
					ID: id.(string),
				})
			}
		}
	}

	if len(credentials) == 0 {
		return nil
	}
	return credentials
}

func flattenCredentials(creds []common.CommonIDName) []interface{} {
	if len(creds) == 0 {
		return nil
	}

	idList := make([]string, len(creds))
	for i, c := range creds {
		idList[i] = c.ID
	}

	// Ensure it's a TypeSet inside the block
	m := map[string]interface{}{
		"id": schema.NewSet(schema.HashString, stringSliceToInterfaceSlice(idList)),
	}

	return []interface{}{m}
}

func stringSliceToInterfaceSlice(input []string) []interface{} {
	out := make([]interface{}, len(input))
	for i, v := range input {
		out[i] = v
	}
	return out
}
