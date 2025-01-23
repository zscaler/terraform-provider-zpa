package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimattributeheader"
)

func dataSourceScimAttributeHeader() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScimAttributeHeaderRead,
		Schema: map[string]*schema.Schema{
			"canonical_values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"case_sensitive": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"idp_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"multivalued": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mutability": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"returned": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uniqueness": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceScimAttributeHeaderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Ensure the service is set correctly before proceeding
	if service == nil {
		return diag.FromErr(fmt.Errorf("ScimAttributeHeader service is not available"))
	}

	// Prepare variables to hold the response objects
	var resp *scimattributeheader.ScimAttributeHeader
	var idpResp *idpcontroller.IdpController

	// Extract IDP-related attributes
	idpId, okidpId := d.Get("idp_id").(string)
	idpName, okIdpName := d.Get("idp_name").(string)

	// Check for presence of either IDP name or ID
	if (!okIdpName && !okidpId) || (idpId == "" && idpName == "") {
		log.Printf("[INFO] IDP name or ID is required\n")
		return diag.FromErr(fmt.Errorf("IDP name or ID is required"))
	}

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

	// Retrieve SCIM attribute header by ID or name
	id, idExists := d.Get("id").(string)
	name, nameExists := d.Get("name").(string)

	if idExists && id != "" {
		// Fetch by ID
		res, _, err := scimattributeheader.Get(ctx, service, idpResp.ID, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	} else if nameExists && name != "" {
		// Fetch by name
		res, _, err := scimattributeheader.GetByName(ctx, service, name, idpResp.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	// Populate the resource data if the response is not nil
	if resp != nil {
		values, err := scimattributeheader.GetValues(ctx, service, resp.IdpID, resp.ID)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resp.ID)
		_ = d.Set("canonical_values", resp.CanonicalValues)
		_ = d.Set("case_sensitive", resp.CaseSensitive)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("data_type", resp.DataType)
		_ = d.Set("description", resp.Description)
		_ = d.Set("idp_id", resp.IdpID)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("multivalued", resp.MultiValued)
		_ = d.Set("mutability", resp.Mutability)
		_ = d.Set("name", resp.Name)
		_ = d.Set("required", resp.Required)
		_ = d.Set("returned", resp.Returned)
		_ = d.Set("schema_uri", resp.SchemaURI)
		_ = d.Set("uniqueness", resp.Uniqueness)
		_ = d.Set("values", values)
	} else {
		return diag.FromErr(fmt.Errorf("no SCIM attribute with name '%s' and IDP name '%s', or ID '%s' was found", name, idpName, id))
	}
	return nil
}
