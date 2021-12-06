---
subcategory: "Provisioning Key"
layout: "zpa"
page_title: "ZPA: provisioning_key"
description: |-
  Creates a ZPA Provisioning Key for Service Edge and/or App Connector Groups.

---

# zpa_provisioning_key (Resource)

The **zpa_provisioning_key** resource provides creates a provisioning key in the Zscaler Private Access portal. This resource can then be referenced in the following ZPA resources:

1. App Connector Groups
2. Service Edge Groups

## Example Usage

```hcl
# Create Provisioning Key for Service Edge Group
resource "zpa_provisioning_key" "usa_provisioning_key" {
  name                  = "AWS Provisioning Key"
  association_type      = "SERVICE_EDGE_GRP"
  max_usage             = "10"
  enrollment_cert_id    = data.zpa_enrollment_cert.service_edge.id
  zcomponent_id         = zpa_service_edge_group.service_edge_group_nyc.id
}

// Create a Service Edge Group
resource "zpa_service_edge_group" "service_edge_group_nyc" {
  name                  = "Service Edge Group New York"
  description           = "Service Edge Group in New York"
  upgrade_day           = "SUNDAY"
  upgrade_time_in_secs  = "66600"
  latitude              = "40.7128"
  longitude             = "-73.935242"
  location              = "New York, NY, USA"
  version_profile_id    = "0"
}

// Retrieve the Service Edge Enrollment Certificate
data "zpa_enrollment_cert" "service_edge" {
    name = "Service Edge"
}
```

```hcl
// Create Provisioning Key for App Connector Group
resource "zpa_provisioning_key" "canada_provisioning_key" {
  name                  = "Canada Provisioning Key"
  association_type      = "CONNECTOR_GRP"
  max_usage             = "10"
  enrollment_cert_id    = data.zpa_enrollment_cert.connector.id
  zcomponent_id         = zpa_app_connector_group.canada_connector_group.id
}

// Create an App Connector Group
resource "zpa_app_connector_group" "canada_connector_group" {
  name                          = "Canada Connector Group"
  description                   = "Canada Connector Group"
  enabled                       = true
  city_country                  = "Toronto, CA"
  country_code                  = "CA"
  latitude                      = "43.6532"
  longitude                     = "79.3832"
  location                      = "Toronto, ON, Canada"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  dns_query_type                = "IPV4"
}

// Retrieve the Connector Enrollment Certificate
data "zpa_enrollment_cert" "connector" {
    name = "Connector"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the provisioning key.
* `max_usage` - (Required) The maximum number of instances where this provisioning key can be used for enrolling an App Connector or Service Edge.
* `enrollment_cert_id` - (Required) ID of the enrollment certificate that can be used for this provisioning key. `ID` of the existing enrollment certificate that has the private key
* `zcomponentId` - (Required) ID of the existing App Connector or Service Edge Group.
* `association_type` (Required) Specifies the provisioning key type for App Connectors or ZPA Private Service Edges. The supported values are `CONNECTOR_GRP` and `SERVICE_EDGE_GRP`

## Import

Provisioning key can be imported by using `<PROVISIONING KEY ID>` or `<PROVISIONING KEY NAME>` as the import ID.

For example:

```shell
terraform import zpa_provisioning_key.example <provisioning_key_id>
```

or

```shell
terraform import zpa_provisioning_key.example <provisioning_key_name>
```
