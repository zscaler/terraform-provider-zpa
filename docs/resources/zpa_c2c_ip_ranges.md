---
page_title: "zpa_c2c_ip_ranges Resource - terraform-provider-zpa"
subcategory: "C2C IP Ranges"
description: |-
  Official documentation https://help.zscaler.com/zpa/adding-ip-ranges
  API documentation https://help.zscaler.com/zpa/adding-ip-ranges
  Creates and manages ZPA C2C IP Ranges in Zscaler Private Access cloud.
---

# zpa_c2c_ip_ranges (Resource)

* [Official documentation](https://help.zscaler.com/zpa/adding-ip-ranges)
* [API documentation](https://help.zscaler.com/zpa/adding-ip-ranges)

The **zpa_c2c_ip_ranges** resource creates a C2C IP Ranges in the Zscaler Private Access cloud.

## Example Usage - Using IP Range

```terraform
resource "zpa_c2c_ip_ranges" "this" {
  name            = "Terraform_IP_Range01"
  description     = "Terraform_IP_Range01"
  enabled         = true
  location_hint   = "Created_via_Terraform"
  ip_range_begin  = "192.168.1.1"
  ip_range_end    = "192.168.1.254"
  location        = "San Jose, CA, USA"
  sccm_flag       = true
  country_code    = "US"
  latitude_in_db  = "37.33874"
  longitude_in_db = "-121.8852525"
}
```

## Example Usage - Using Subnet CIDR

```terraform
resource "zpa_c2c_ip_ranges" "this" {
  name            = "Terraform_IP_Range01"
  description     = "Terraform_IP_Range01"
  enabled         = true
  location_hint   = "Created_via_Terraform"
  subnet_cidr     = "192.168.1.0/24"
  location        = "San Jose, CA, USA"
  sccm_flag       = true
  country_code    = "US"
  latitude_in_db  = "37.33874"
  longitude_in_db = "-121.8852525"
}
```

## Schema

### Required

- `name` (String) - Name of the C2C IP Ranges

### Optional

- `id` (String) - The ID of the C2C IP Ranges
- `description` (String) - Description of the C2C IP Ranges
- `enabled` (Boolean) - Whether the C2C IP Ranges is enabled
- `ip_range_begin` (String) - Beginning IP address of the range
- `ip_range_end` (String) - Ending IP address of the range
- `location` (String) - Location of the C2C IP Ranges
- `location_hint` (String) - Location hint for the C2C IP Ranges
- `sccm_flag` (Boolean) - SCCM flag for the C2C IP Ranges
- `subnet_cidr` (String) - Subnet CIDR for the C2C IP Ranges
- `country_code` (String) - Country code for the C2C IP Ranges
- `latitude_in_db` (String) - Latitude in database for the C2C IP Ranges
- `longitude_in_db` (String) - Longitude in database for the C2C IP Ranges

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/SecurityGeekIO/zscaler-terraformer)

C2C IP Ranges can be imported by using `<RANGE ID>` or `<RANGE NAME>` as the import ID

For example:

```shell
terraform import zpa_c2c_ip_ranges.example <range_id>
```

or

```shell
terraform import zpa_c2c_ip_ranges.example <range_name>
```
