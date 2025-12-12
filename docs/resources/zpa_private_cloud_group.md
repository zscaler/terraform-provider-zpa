---
page_title: "zpa_private_cloud_group Resource - terraform-provider-zpa"
subcategory: "Private Cloud Group"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-private-cloud-controller-groups
  API documentation https://help.zscaler.com/zpa/about-private-cloud-controller-groups
  Creates and manages ZPA Private Cloud Group in Zscaler Private Access cloud.
---

# zpa_private_cloud_group (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-private-cloud-controller-groups)
* [API documentation](https://help.zscaler.com/zpa/about-private-cloud-controller-groups)

The **zpa_private_cloud_group** resource creates a private cloud group in the Zscaler Private Access cloud.

## Example Usage

```terraform
resource "zpa_private_cloud_group" "this" {
  name                     = "PrivateCloudGroup01"
  description              = "Example private cloud group"
  enabled                  = true
  city_country             = "San Jose, US"
  latitude                 = "37.33874"
  longitude                = "-121.8852525"
  location                 = "San Jose, CA, USA"
  upgrade_day              = "SUNDAY"
  upgrade_time_in_secs     = "66600"
  site_id                  = "72058304855088543"
  version_profile_id       = "0"
  override_version_profile = true
  is_public                = "TRUE"
}
```

## Schema

### Required

- `name` (String) - Name of the Private Cloud Group

### Optional

- `id` (String) - The ID of the Private Cloud Group
- `city_country` (String) - City and country of the Private Cloud Group
- `country_code` (String) - Country code of the Private Cloud Group
- `description` (String) - Description of the Private Cloud Group
- `enabled` (Boolean) - Whether this Private Cloud Group is enabled or not
- `is_public` (String) - Whether the Private Cloud Group is public
- `latitude` (String) - Latitude of the Private Cloud Group. Integer or decimal. With values in the range of -90 to 90
- `location` (String) - Location of the Private Cloud Group
- `longitude` (String) - Longitude of the Private Cloud Group. Integer or decimal. With values in the range of -180 to 180
- `override_version_profile` (Boolean) - Whether the default version profile of the Private Cloud Group is applied or overridden
- `microtenant_id` (String) - Microtenant ID for the Private Cloud Group
- `site_id` (String) - Site ID for the Private Cloud Group
- `upgrade_day` (String) - Private Cloud Controllers in this group will attempt to update to a newer version of the software during this specified day. Supported values: `SUNDAY`, `MONDAY`, `TUESDAY`, `WEDNESDAY`, `THURSDAY`, `FRIDAY`, `SATURDAY`
- `upgrade_time_in_secs` (String) - Private Cloud Controllers in this group will attempt to update to a newer version of the software during this specified time. Integer in seconds (i.e., -66600). The integer should be greater than or equal to 0 and less than 86400, in 15 minute intervals
- `version_profile_id` (String) - ID of the version profile for the Private Cloud Group

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/SecurityGeekIO/zscaler-terraformer)

Private Cloud Group can be imported by using `<GROUP ID>` or `<GROUP NAME>` as the import ID.

```shell
terraform import zpa_private_cloud_group.example <group_id>
```

or

```shell
terraform import zpa_private_cloud_group.example <group_name>
```
