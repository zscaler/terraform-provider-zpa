---
page_title: "zpa_c2c_ip_ranges Data Source - terraform-provider-zpa"
subcategory: "C2C IP Ranges"
description: |-
  Official documentation https://help.zscaler.com/zpa/adding-ip-ranges
  API documentation https://help.zscaler.com/zpa/adding-ip-ranges
  Get information about ZPA C2C IP Ranges in Zscaler Private Access cloud.
---

# zpa_c2c_ip_ranges (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/adding-ip-ranges)
* [API documentation](https://help.zscaler.com/zpa/adding-ip-ranges)

The **zpa_c2c_ip_ranges** data source to get information about a C2C IP range created in the Zscaler Private Access cloud.

## Example Usage

```terraform
# ZPA C2C IP Ranges Data Source by Name
data "zpa_c2c_ip_ranges" "this" {
 name = "Range01"
}
```

```terraform
# ZPA C2C IP Ranges Data Source by ID
data "zpa_c2c_ip_ranges" "this" {
 id = "1234567890"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the C2C IP range to be exported.

### Optional

* `id` - (Optional) The ID of the C2C IP range to be exported.

### Read-Only

* `available_ips` - (String) Available IPs in the C2C IP range
* `country_code` - (String) Country code of the C2C IP range
* `creation_time` - (String) Creation time of the C2C IP range
* `customer_id` - (String) Customer ID of the C2C IP range
* `description` - (String) Description of the C2C IP range
* `enabled` - (Boolean) Whether the C2C IP range is enabled
* `ip_range_begin` - (String) Beginning IP address of the range
* `ip_range_end` - (String) Ending IP address of the range
* `is_deleted` - (String) Whether the C2C IP range is deleted
* `latitude_in_db` - (String) Latitude in database for the C2C IP range
* `location` - (String) Location of the C2C IP range
* `location_hint` - (String) Location hint for the C2C IP range
* `longitude_in_db` - (String) Longitude in database for the C2C IP range
* `modified_by` - (String) Modified by information for the C2C IP range
* `modified_time` - (String) Modified time of the C2C IP range
* `sccm_flag` - (Boolean) SCCM flag for the C2C IP range
* `subnet_cidr` - (String) Subnet CIDR for the C2C IP range
* `total_ips` - (String) Total IPs in the C2C IP range
* `used_ips` - (String) Used IPs in the C2C IP range
