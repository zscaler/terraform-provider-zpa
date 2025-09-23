package zpa

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/c2c_ip_ranges"
)

func resourceC2CIPRanges() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceC2CIPRangesCreate,
		ReadContext:   resourceC2CIPRangesRead,
		UpdateContext: resourceC2CIPRangesUpdate,
		DeleteContext: resourceC2CIPRangesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.Service

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
					resp, _, err := c2c_ip_ranges.GetByName(ctx, service, id)
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the C2C IP Ranges",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the C2C IP Ranges",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the C2C IP Ranges is enabled",
			},
			"ip_range_begin": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Beginning IP address of the range",
			},
			"ip_range_end": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Ending IP address of the range",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location of the C2C IP Ranges",
			},
			"location_hint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location hint for the C2C IP Ranges",
			},
			"sccm_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "SCCM flag for the C2C IP Ranges",
			},
			"subnet_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet CIDR for the C2C IP Ranges",
			},
			"country_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Country code for the C2C IP Ranges",
			},
			"latitude_in_db": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Latitude in database for the C2C IP Ranges",
			},
			"longitude_in_db": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Longitude in database for the C2C IP Ranges",
			},
		},
	}
}

func resourceC2CIPRangesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandC2CIPRanges(d)
	log.Printf("[INFO] Creating C2C IP Ranges with request:\n%+v\n", req)
	resp, _, err := c2c_ip_ranges.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created C2C IP Ranges. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceC2CIPRangesRead(ctx, d, meta)
}

func resourceC2CIPRangesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id := d.Id()
	log.Printf("[INFO] Getting C2C IP Ranges with id: %v\n", id)
	resp, _, err := c2c_ip_ranges.Get(ctx, service, id)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Got C2C IP Ranges:\n%+v\n", resp)

	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("ip_range_begin", resp.IpRangeBegin)
	_ = d.Set("ip_range_end", resp.IpRangeEnd)
	_ = d.Set("location", resp.Location)
	_ = d.Set("location_hint", resp.LocationHint)
	_ = d.Set("sccm_flag", resp.SccmFlag)
	_ = d.Set("subnet_cidr", resp.SubnetCidr)
	_ = d.Set("country_code", resp.CountryCode)
	_ = d.Set("latitude_in_db", resp.LatitudeInDb)
	_ = d.Set("longitude_in_db", resp.LongitudeInDb)

	return nil
}

func resourceC2CIPRangesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id := d.Id()
	log.Printf("[INFO] Updating C2C IP Ranges with id: %v\n", id)
	req := expandC2CIPRanges(d)
	if _, err := c2c_ip_ranges.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceC2CIPRangesRead(ctx, d, meta)
}

func resourceC2CIPRangesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id := d.Id()
	log.Printf("[INFO] Deleting C2C IP Ranges with id: %v\n", id)
	if _, err := c2c_ip_ranges.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Deleted C2C IP Ranges with id: %v\n", id)

	return nil
}

func expandC2CIPRanges(d *schema.ResourceData) c2c_ip_ranges.IPRanges {
	return c2c_ip_ranges.IPRanges{
		ID:            d.Get("id").(string),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Enabled:       d.Get("enabled").(bool),
		IpRangeBegin:  d.Get("ip_range_begin").(string),
		IpRangeEnd:    d.Get("ip_range_end").(string),
		Location:      d.Get("location").(string),
		LocationHint:  d.Get("location_hint").(string),
		SccmFlag:      d.Get("sccm_flag").(bool),
		SubnetCidr:    d.Get("subnet_cidr").(string),
		CountryCode:   d.Get("country_code").(string),
		LatitudeInDb:  d.Get("latitude_in_db").(string),
		LongitudeInDb: d.Get("longitude_in_db").(string),
	}
}
