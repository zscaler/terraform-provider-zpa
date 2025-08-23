package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/c2c_ip_ranges"
)

func dataSourceC2CIPRanges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceC2CIPRangesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"available_ips": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"country_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ip_range_begin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_range_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latitude_in_db": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_hint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"longitude_in_db": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sccm_flag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"subnet_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_ips": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_ips": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceC2CIPRangesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *c2c_ip_ranges.IPRanges
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for C2C IP ranges %s\n", id)
		res, _, err := c2c_ip_ranges.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for C2C IP ranges name %s\n", name)
		// Since there's no GetByName function, we'll get all and filter by name
		allRanges, _, err := c2c_ip_ranges.GetAll(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, ipRange := range allRanges {
			if ipRange.Name == name {
				resp = ipRange
				break
			}
		}
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("available_ips", resp.AvailableIps)
		_ = d.Set("country_code", resp.CountryCode)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("customer_id", resp.CustomerId)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("ip_range_begin", resp.IpRangeBegin)
		_ = d.Set("ip_range_end", resp.IpRangeEnd)
		_ = d.Set("is_deleted", resp.IsDeleted)
		_ = d.Set("latitude_in_db", resp.LatitudeInDb)
		_ = d.Set("location", resp.Location)
		_ = d.Set("location_hint", resp.LocationHint)
		_ = d.Set("longitude_in_db", resp.LongitudeInDb)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("sccm_flag", resp.SccmFlag)
		_ = d.Set("subnet_cidr", resp.SubnetCidr)
		_ = d.Set("total_ips", resp.TotalIps)
		_ = d.Set("used_ips", resp.UsedIps)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any C2C IP ranges with name '%s' or id '%s'", name, id))
	}

	return nil
}
