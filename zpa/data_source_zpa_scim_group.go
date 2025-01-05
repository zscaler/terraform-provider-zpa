package zpa

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimgroup"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceScimGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScimGroupRead,
		Schema: map[string]*schema.Schema{
			"creation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"idp_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"idp_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"idp_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceScimGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Ensure the service is set correctly before proceeding
	if service == nil {
		return diag.FromErr(fmt.Errorf("ScimGroup service is not available"))
	}

	var resp *scimgroup.ScimGroup
	idpId, okidpId := d.Get("idp_id").(string)
	idpName, okIdpName := d.Get("idp_name").(string)

	// Ensure either IDP name or ID is provided
	if (!okIdpName && !okidpId) || (idpId == "" && idpName == "") {
		log.Printf("[INFO] IDP name or ID is required\n")
		return diag.FromErr(fmt.Errorf("IDP name or ID is required"))
	}

	var idpResp *idpcontroller.IdpController
	var err error
	// Get IDP Controller by ID or name
	if idpId != "" {
		idpResp, _, err = idpcontroller.Get(ctx, service, idpId)
		if err != nil || idpResp == nil {
			log.Printf("[INFO] Couldn't find IDP by ID: %s\n", idpId)
			return diag.FromErr(fmt.Errorf("error fetching IDP by ID: %w", err))
		}
	} else {
		idpResp, _, err = idpcontroller.GetByName(ctx, service, idpName)
		if err != nil || idpResp == nil {
			log.Printf("[INFO] Couldn't find IDP by name: %s\n", idpName)
			return diag.FromErr(fmt.Errorf("error fetching IDP by name: %w", err))
		}
	}

	// Declare variables for id and name to ensure they are accessible in the final error message
	id, idExists := d.Get("id").(string)
	name, nameExists := d.Get("name").(string)

	// Retrieve SCIM group by ID or name
	if idExists && id != "" {
		res, _, err := scimgroup.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	} else if nameExists && name != "" && idpResp != nil {
		// Check idpResp is non-nil before accessing its fields
		res, _, err := scimgroup.GetByName(ctx, service, name, idpResp.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	// Set resource data if the response is not nil
	if resp != nil {
		d.SetId(strconv.FormatInt(int64(resp.ID), 10))
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("idp_group_id", resp.IdpGroupID)
		_ = d.Set("idp_id", resp.IdpID)
		_ = d.Set("idp_name", resp.IdpName)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
	} else {
		return diag.FromErr(fmt.Errorf("no SCIM group with name '%s' and IDP name '%s', or ID '%s' was found", name, idpName, id))
	}
	return nil
}
