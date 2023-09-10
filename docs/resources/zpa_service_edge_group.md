---
subcategory: "Service Edge Group"
layout: "zscaler"
page_title: "ZPA: service_edge_group"
description: |-
  Creates and manages ZPA Service Edge Group details.
---

# Resource: zpa_service_edge_group

The **zpa_service_edge_group** resource creates a service edge group in the Zscaler Private Access cloud. This resource can then be referenced in a service edge connector.

## Example Usage

```hcl
# ZPA Service Edge Group resource - Trusted Network
resource "zpa_service_edge_group" "service_edge_group_sjc" {
  name                 = "Service Edge Group San Jose"
  description          = "Service Edge Group in San Jose"
  enabled              = true
  is_public            = true
  upgrade_day          = "SUNDAY"
  upgrade_time_in_secs = "66600"
  latitude             = "37.3382082"
  longitude            = "-121.8863286"
  location             = "San Jose, CA, USA"
  version_profile_name = "New Release"
  trusted_networks {
    id = [ data.zpa_trusted_network.example.id ]
  }
}
```

```hcl
# ZPA Service Edge Group resource - No Trusted Network
resource "zpa_service_edge_group" "service_edge_group_nyc" {
  name                 = "Service Edge Group New York"
  description          = "Service Edge Group in New York"
  enabled              = true
  is_public            = true
  upgrade_day          = "SUNDAY"
  upgrade_time_in_secs = "66600"
  latitude             = "40.7128"
  longitude            = "-73.935242"
  location             = "New York, NY, USA"
  version_profile_name = "New Release"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Service Edge Group.
* `latitude` - (Required) Latitude for the Service Edge Group. Integer or decimal with values in the range of `-90` to `90`
* `longitude` - (Required) Longitude for the Service Edge Group. Integer or decimal with values in the range of `-180` to `180`
* `location` - (Required) Location for the Service Edge Group.
* `description` - (Optional) Description of the Service Edge Group.
* `enabled` - (Optional) Whether this Service Edge Group is enabled or not. Default value: `true` Supported values: `true`, `false`
* `city_country` - (Optional) This field controls dynamic discovery of the servers.
* `country_code` - (Optional) This field is an array of app-connector-id only.
* `is_public` - (Optional) Enable or disable public access for the Service Edge Group. Default value: `false` Supported values: `true`, `false`

* `override_version_profile` - (Optional) Whether the default version profile of the App Connector Group is applied or overridden. Default: `false` Supported values: `true`, `false`
* `version_profile_id` - (Optional) ID of the version profile. To learn more, see Version Profile Use Cases. Supported values are:
  * ``0`` = ``Default``
  * ``1`` = ``Previous Default``
  * ``2`` = ``New Release``
* `version_profile_name` - (Optional)
  * ``Default`` = ``0``
  * ``Previous Default`` = ``1``
  * ``New Release`` = ``2``
* `service_edges` - (Optional)
* `trusted_networks` - (Optional) Trusted networks for this Service Edge Group. List of trusted network objects
* `upgrade_day` - (Optional) Service Edges in this group will attempt to update to a newer version of the software during this specified day. Default value: `SUNDAY` List of valid days (i.e., Sunday, Monday)
* `upgrade_time_in_secs` - (Optional) Service Edges in this group will attempt to update to a newer version of the software during this specified time. Default value: `66600` Integer in seconds (i..e, 66600). The integer must be greater than or equal to 0 and less than `86400`, in `15` minute intervals
* `microtenant_id` (Optional) The ID of the microtenant the resource is to be associated with.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Service Edge Group can be imported; use `<SERVER EDGE GROUP ID>` or `<SERVER EDGE GROUP NAME>` as the import ID.

For example:

```shell
terraform import zpa_service_edge_group.example <service_edge_group_id>
```

or

```shell
terraform import zpa_service_edge_group.example <service_edge_group_name>
```
