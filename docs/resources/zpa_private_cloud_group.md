---
page_title: "zpa_private_cloud_group Resource - terraform-provider-zpa"
subcategory: "Private Cloud Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-private-cloud-controller-groups
  API documentation https://help.zscaler.com/zpa/about-private-cloud-controller-groups
  Creates and manages ZPA Private Cloud Group in Zscaler Private Access cloud.
---

# zpa_private_cloud_group (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-private-cloud-controller-groups)
* [API documentation](https://help.zscaler.com/zpa/about-private-cloud-controller-groups)

The **zpa_private_cloud_group** resource creates a private cloud group in the Zscaler Private Access cloud.

## Private Cloud Connector Onboarding Methods

Private Cloud Connectors can be onboarded into ZPA in two ways. This resource supports both:

1. **OAuth2 user codes** *(recommended for new deployments)* - Set `user_codes` with the codes generated on each Private Cloud Connector VM. The provider creates the group and then calls the OAuth2 user code verification API to enroll the connectors.
2. **Provisioning key** *(legacy / still supported)* - Create the group with this resource, then create a `zpa_provisioning_key` referencing it. The key is then injected into the Private Cloud Connector VM at deployment time.

In **both** methods, the Private Cloud Connector enrollment requires an `enrollment_cert_id`. You can either:
- Set `enrollment_cert_id` explicitly using the `zpa_enrollment_cert` data source, or
- Omit it entirely - the provider will automatically look up the **"Connector"** enrollment certificate by name and populate the ID for you.

## Example Usage - OAuth2 enrollment with user codes (Explicit Enrollment Certificate)

Set the enrollment certificate explicitly and provide the user codes displayed on the Private Cloud Connector VMs after deployment. The provider will create the group and then call the user code verification API to complete enrollment.

```terraform
data "zpa_enrollment_cert" "connector" {
  name = "Connector"
}

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
  enrollment_cert_id = data.zpa_enrollment_cert.connector.id
  user_codes         = ["CODE_FROM_VM_1", "CODE_FROM_VM_2"]
}
```

## Example Usage - OAuth2 enrollment with user codes (Auto-resolved Enrollment Certificate)

Omit `enrollment_cert_id` entirely and the provider will automatically resolve the **"Connector"** enrollment certificate for you. This is the simplest configuration and is functionally equivalent to the explicit example above.

```terraform
resource "zpa_private_cloud_group" "example" {
  name                     = "PrivateCloudGroup01"
  description              = "Example private cloud group"
  enabled                  = true
  country_code             = "US"
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

  user_codes = ["CODE_FROM_VM_1", "CODE_FROM_VM_2"]
}
```

## Example Usage - Enrolling Private Cloud Connectors Via Provisioning Key (Explicit Enrollment Certificate)

Create the Private Cloud Connector Group, then create a `zpa_provisioning_key` that references the group's ID. The provisioning key is then injected into the Private Cloud Connector VM at deployment time.

```terraform
data "zpa_enrollment_cert" "connector" {
  name = "Connector"
}

resource "zpa_private_cloud_group" "example" {
  name                 = "Example"
  description          = "Example"
  enabled              = true
  city_country         = "San Jose, CA"
  country_code         = "US"
  latitude             = "37.338"
  longitude            = "-121.8863"
  location             = "San Jose, CA, US"
  upgrade_day          = "SUNDAY"
  upgrade_time_in_secs = "66600"
  dns_query_type       = "IPV4_IPV6"

  enrollment_cert_id = data.zpa_enrollment_cert.connector.id
}

resource "zpa_provisioning_key" "example" {
  name               = "ProvisioningKey01"
  association_type   = "CONNECTOR_GRP"
  max_usage          = "10"
  enrollment_cert_id = data.zpa_enrollment_cert.connector.id
  zcomponent_id      = zpa_private_cloud_group.example.id
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
[Visit](https://github.com/zscaler/zscaler-terraformer)

Private Cloud Group can be imported by using `<GROUP ID>` or `<GROUP NAME>` as the import ID.

```shell
terraform import zpa_private_cloud_group.example <group_id>
```

or

```shell
terraform import zpa_private_cloud_group.example <group_name>
```
