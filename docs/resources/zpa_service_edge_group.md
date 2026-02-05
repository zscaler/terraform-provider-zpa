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

## Example Usage - Using Version Profile Name

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
  version_profile_id   = data.zpa_customer_version_profile.this.id 
}
```

## Example Usage - Using Version Profile ID

data "zpa_customer_version_profile" "this" {
  name = "New Release"
}

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

## Example Usage - OAuth2 enrollment with user codes

When enrolling Service Edges via OAuth2, set the enrollment certificate and provide the user codes displayed on the Service Edge VMs after deployment. The provider will create the group and then call the user code verification API to complete enrollment.

```terraform
data "zpa_enrollment_cert" "service_edge" {
  name = "Service Edge"
}

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

  enrollment_cert_id   = data.zpa_enrollment_cert.service_edge.id
  user_codes           = ["CODE_FROM_VM_1", "CODE_FROM_VM_2"]

  trusted_networks {
    id = [data.zpa_trusted_network.example.id]
  }
}
```

## Schema

### Required

The following arguments are supported:

- `name` - (String) Name of the Service Edge Group.
- `latitude` - (String) Latitude for the Service Edge Group. Integer or decimal with values in the range of `-90` to `90`
- `longitude` - (String) Longitude for the Service Edge Group. Integer or decimal with values in the range of `-180` to `180`
- `location` - (String) Location of the App Connector Group. i.e ``"San Jose, CA, USA"``
- `city_country` - (String) The city and country of the App Connector i.e ``"San Jose, US"``
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

- `version_profile_id` - (String) The unique identifier of the version profile. Supported values are:
  - ``0`` = ``Default``
  - ``1`` = ``Previous Default``
  - ``2`` = ``New Release``

  **NOTE:** In order to retrieve other version profile IDs, you can leverage the data source `zpa_customer_version_profile`

- `version_profile_name` - (String) The unique identifier of the version profile. Supported values are:
  - ``Default``, ``Previous Default``, ``New Release``, ``Default - el8``, ``New Release - el8``, ``Previous Default - el8``

- `upgrade_day` - (Strings) Service Edges in this group will attempt to update to a newer version of the software during this specified day. Default value: `SUNDAY` List of valid days (i.e., Sunday, Monday)
- `upgrade_time_in_secs` - (Strings) Service Edges in this group will attempt to update to a newer version of the software during this specified time. Default value: `66600` Integer in seconds (i..e, 66600). The integer must be greater than or equal to 0 and less than `86400`, in `15` minute intervals
- `use_in_dr_mode` - (Boolean) Whether or not the App Connector Group is designated for disaster recovery. Supported values: `true`, `false`
- `microtenant_id` (Strings) The ID of the microtenant the resource is to be associated with.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

### OAuth2 enrollment (optional)

- `enrollment_cert_id` - (String) ID of the enrollment certificate used for OAuth2 enrollment. When set along with `user_codes`, the provider will enroll Service Edges via the OAuth2 user code verification API. Use the data source `zpa_enrollment_cert` to look up the certificate (e.g. name = "Service Edge").
- `user_codes` - (Set of String) User codes from deployed Service Edge VMs for OAuth2 enrollment. When provided together with `enrollment_cert_id`, the provider calls the user code verification API to enroll the service edges. Obtain these codes from the Service Edge VM after deployment (they are displayed during the OAuth2 enrollment flow).

- `trusted_networks` - (Block Set) Trusted networks for this Service Edge Group. List of trusted network objects Maximum 1 block allowed.
    - `id` - (List of Strings) The unique identifier of the trusted network.

- `service_edges` - (Block Set) The list of ZPA Private Service Edges in the ZPA Private Service Edge Group. Maximum 1 block allowed.
    - `id` - (List of Strings) The unique identifier of the ZPA Private Service Edge.

  ### New `service_edges` Behavior
  - **When omitted**: Terraform will ignore service edge membership completely (no drift detection)
  - **When specified**: Terraform will enforce exact membership matching
    - You must include all required service edge IDs in the list
    - Any discrepancy between configuration and actual state will be reported as drift

  ### Important Notes
  ⚠️ **Deprecation Notice**: The `service_edges` block will be deprecated in a future release  
  🔧 **External Management**: Service edge membership is typically managed outside Terraform  
  💡 **Recommendation**: Only use this block if you require Terraform to explicitly manage membership

  ### Migration Guidance
  If you're currently using this block but don't need strict membership control:
  1. Remove the `service_edges` block from your configuration
  2. Run `terraform apply` to update the state

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
