---
page_title: "zpa_private_cloud_group Data Source - terraform-provider-zpa"
subcategory: "Private Cloud Group"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-private-cloud-controller-groups
  API documentation https://help.zscaler.com/zpa/about-private-cloud-controller-groups
  Get information about ZPA Private Cloud Group in Zscaler Private Access cloud.
---

# zpa_private_cloud_group (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-private-cloud-controller-groups)
* [API documentation](https://help.zscaler.com/zpa/about-private-cloud-controller-groups)

The **zpa_private_cloud_group** data source to get information about a private cloud group in the Zscaler Private Access cloud.

## Example Usage

```terraform
# ZPA Private Cloud Group Data Source
data "zpa_private_cloud_group" "foo" {
  name = "DataCenter"
}
```

```terraform
# ZPA Private Cloud Group Data Source
data "zpa_private_cloud_group" "foo" {
  id = "123456789"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the private cloud group to be exported.

### Optional

* `id` - (Optional) The ID of the private cloud group to be exported.
* `microtenant_id` - (Optional) Microtenant ID for the private cloud group.

### Read-Only

* `city_country` - (String) City and country of the Private Cloud Group
* `country_code` - (String) Country code of the Private Cloud Group
* `description` - (String) Description of the Private Cloud Group
* `enabled` - (Boolean) Whether this Private Cloud Group is enabled or not
* `geo_location_id` - (String) Geo location ID of the Private Cloud Group
* `is_public` - (String) Whether the Private Cloud Group is public
* `latitude` - (String) Latitude of the Private Cloud Group
* `location` - (String) Location of the Private Cloud Group
* `longitude` - (String) Longitude of the Private Cloud Group
* `override_version_profile` - (Boolean) Whether the default version profile of the Private Cloud Group is applied or overridden
* `read_only` - (Boolean) Whether the Private Cloud Group is read-only
* `restriction_type` - (String) Restriction type of the Private Cloud Group
* `microtenant_name` - (String) Microtenant name for the Private Cloud Group
* `site_id` - (String) Site ID for the Private Cloud Group
* `site_name` - (String) Site name for the Private Cloud Group
* `upgrade_day` - (String) Private Cloud Controllers in this group will attempt to update to a newer version of the software during this specified day
* `upgrade_time_in_secs` - (String) Private Cloud Controllers in this group will attempt to update to a newer version of the software during this specified time
* `version_profile_id` - (String) ID of the version profile for the Private Cloud Group
* `zscaler_managed` - (Boolean) Whether the Private Cloud Group is managed by Zscaler
