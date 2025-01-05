package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/samlattribute"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSamlAttribute() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSamlAttributeRead,
		Schema: map[string]*schema.Schema{
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"idp_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"idp_name": {
				Type:     schema.TypeString,
				Optional: true,
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
				Computed: true,
				Optional: true,
			},
			"saml_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_attribute": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSamlAttributeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *samlattribute.SamlAttribute
	idpId, okidpId := d.Get("idp_id").(string)
	idpName, okIdpName := d.Get("idp_name").(string)

	// Check that either `idp_id` or `idp_name` is provided
	if (!okIdpName && !okidpId) || (idpId == "" && idpName == "") {
		log.Printf("[INFO] IDP name or ID is required\n")
		return diag.FromErr(fmt.Errorf("idp name or id is required"))
	}

	var idpResp *idpcontroller.IdpController
	// Fetch the IDP Controller by ID or name
	if idpId != "" {
		resp, _, err := idpcontroller.Get(ctx, service, idpId)
		if err != nil || resp == nil {
			log.Printf("[INFO] Couldn't find IDP by ID: %s\n", idpId)
			return diag.FromErr(fmt.Errorf("error fetching IDP by ID: %w", err))
		}
		idpResp = resp
	} else {
		resp, _, err := idpcontroller.GetByName(ctx, service, idpName)
		if err != nil || resp == nil {
			log.Printf("[INFO] Couldn't find IDP by name: %s\n", idpName)
			return diag.FromErr(fmt.Errorf("error fetching IDP by name: %w", err))
		}
		idpResp = resp
	}

	// Declare variables for id and name to ensure they are accessible in the final error message
	id, idExists := d.Get("id").(string)
	name, nameExists := d.Get("name").(string)

	// Retrieve SAML attribute by ID or name
	if idExists && id != "" {
		res, _, err := samlattribute.Get(ctx, service, idpResp.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	} else if nameExists && name != "" {
		res, _, err := samlattribute.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	// Set the resource data if the response is not nil
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("idp_id", resp.IdpID)
		_ = d.Set("idp_name", resp.IdpName)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("saml_name", resp.SamlName)
		_ = d.Set("user_attribute", resp.UserAttribute)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any SAML attribute with name '%s' or id '%s'", name, id))
	}

	return nil
}
