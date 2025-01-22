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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredential"
)

func resourcePRACredentialController() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePRACredentialControllerCreate,
		ReadContext:   resourcePRACredentialControllerRead,
		UpdateContext: resourcePRACredentialControllerUpdate,
		DeleteContext: resourcePRACredentialControllerDelete,
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
					resp, _, err := pracredential.GetByName(ctx, service, id)
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
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the privileged credential",
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
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password that is used to protect the SSH private key. This field is optional",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password associated with the username for the login you want to use for the privileged credential",
			},
			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The SSH private key associated with the username for the login you want to use for the privileged credential",
			},
			"user_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The domain name associated with the username. The domain name only needs to be specified with logging in to an RDP console that is connected to an Active Directory Domain",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: " The username for the login you want to use for the privileged credential",
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

func resourcePRACredentialControllerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandPRACredentialController(d)

	sanitizeFields(&req)
	log.Printf("[INFO] Creating credential controller with request\n%+v\n", req)

	credController, _, err := pracredential.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created credential controller request. ID: %v\n", credController)

	d.SetId(credController.ID)
	return resourcePRACredentialControllerRead(ctx, d, meta)
}

func resourcePRACredentialControllerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := pracredential.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing credential controller %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting credential controller:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("credential_type", resp.CredentialType)
	_ = d.Set("username", resp.UserName)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	return nil
}

func resourcePRACredentialControllerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	log.Printf("[INFO] Updating credential controller ID: %v\n", id)

	req := expandPRACredentialController(d)

	if _, _, err := pracredential.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := pracredential.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourcePRACredentialControllerRead(ctx, d, meta)
}

func resourcePRACredentialControllerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting credential controller ID: %v\n", d.Id())

	if _, err := pracredential.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] credential controller deleted")
	return nil
}

func expandPRACredentialController(d *schema.ResourceData) pracredential.Credential {
	credController := pracredential.Credential{
		ID:             d.Id(),
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		CredentialType: d.Get("credential_type").(string),
		Passphrase:     d.Get("passphrase").(string),
		Password:       d.Get("password").(string),
		PrivateKey:     d.Get("private_key").(string),
		UserDomain:     d.Get("user_domain").(string),
		UserName:       d.Get("username").(string),
		MicroTenantID:  d.Get("microtenant_id").(string),
	}
	return credController
}
