---
page_title: "zpa_service_edge_group Resource - terraform-provider-zpa"
subcategory: "Service Edge Group"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-zpa-private-service-edge-groups
  API documentation https://help.zscaler.com/zpa/configuring-zpa-private-service-edge-groups-using-api
  Creates and manages ZPA Service Edge Group details.
---

# zpa_service_edge_group (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-zpa-private-service-edge-groups)
* [API documentation](https://help.zscaler.com/zpa/configuring-zpa-private-service-edge-groups-using-api)

The **zpa_service_edge_group** resource creates a service edge group in the Zscaler Private Access cloud. This resource can then be referenced in a service edge connector.

## Example Usage

```terraform
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

```terraform
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

## Schema

### Required

The following arguments are supported:

- `name` - (String) Name of the Service Edge Group.
- `latitude` - (String) Latitude for the Service Edge Group. Integer or decimal with values in the range of `-90` to `90`
- `longitude` - (String) Longitude for the Service Edge Group. Integer or decimal with values in the range of `-180` to `180`
- `location` - (String) Location of the App Connector Group. i.e ``"San Jose, CA, USA"``
- `city_country` - (String) Whether Double Encryption is enabled or disabled for the app. i.e ``"San Jose, US"``
- `country_code` - (String) Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

### Optional

In addition to all arguments above, the following attributes are exported:

- `enabled` - (Boolean) Whether this Service Edge Group is enabled or not. Default value: `true` Supported values: `true`, `false`
- `description` - (String) Description of the Service Edge Group.
- `is_public` - (String) Enable or disable public access for the Service Edge Group. Default value: `false` Supported values: `true`, `false`

- `grace_distance_enabled`: Allows ZPA Private Service Edge Groups within the specified distance to be prioritized over a closer ZPA Public Service Edge.
- `grace_distance_value`: Indicates the maximum distance in miles or kilometers to ZPA Private Service Edge groups that would override a ZPA Public Service Edge.
- `grace_distance_value_unit`: Indicates the grace distance unit of measure in miles or kilometers. This value is only required if `grace_distance_enabled` is set to true. Support values are: `MILES` and `KMS`

- `override_version_profile` - (Boolean) Whether the default version profile of the App Connector Group is applied or overridden. Default: `false` Supported values: `true`, `false`
- `version_profile_id` - (String) ID of the version profile. To learn more, see Version Profile Use Cases. Supported values are:
  - ``0`` = ``Default``
  - ``1`` = ``Previous Default``
  - ``2`` = ``New Release``
- `service_edges` - (Block Set) The list of ZPA Private Service Edges in the ZPA Private Service Edge Group.
    - `id` - (List of Strings) The unique identifier of the ZPA Private Service Edge.
- `trusted_networks` - (Block Set) Trusted networks for this Service Edge Group. List of trusted network objects
    - `id` - (List of Strings) The unique identifier of the trusted network.
- `upgrade_day` - (Strings) Service Edges in this group will attempt to update to a newer version of the software during this specified day. Default value: `SUNDAY` List of valid days (i.e., Sunday, Monday)
- `upgrade_time_in_secs` - (Strings) Service Edges in this group will attempt to update to a newer version of the software during this specified time. Default value: `66600` Integer in seconds (i..e, 66600). The integer must be greater than or equal to 0 and less than `86400`, in `15` minute intervals
- `microtenant_id` (Strings) The ID of the microtenant the resource is to be associated with.

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
