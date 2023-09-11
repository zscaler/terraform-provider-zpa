---
subcategory: "Provisioning Key"
layout: "zscaler"
page_title: "ZPA: provisioning_key"
description: |-
  Creates and manages ZPA Provisioning Key for Service Edge and/or App Connector Groups.
---

# Resource: zpa_provisioning_key

The **zpa_provisioning_key** resource provides creates a provisioning key in the Zscaler Private Access portal. This resource can then be referenced in the following ZPA resources:

1. App Connector Groups
2. Service Edge Groups

## Zenith Community - ZPA Provisioning Keys

[![ZPA Terraform provider Video Series Ep3 - Provisioning Keys](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_provisioning_key.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEnCAI/video-zpa-terraform-provider-video-series-ep3-provisioning-keys)

## App Connector Group Provisioning Key Example Usage

```hcl
# Retrieve the Connector Enrollment Certificate
data "zpa_enrollment_cert" "connector" {
    name = "Connector"
}

# Create Provisioning Key for App Connector Group
resource "zpa_provisioning_key" "test_provisioning_key" {
  name                  = test_provisioning_key
  association_type      = "CONNECTOR_GRP"
  max_usage             = "10"
  enrollment_cert_id    = data.zpa_enrollment_cert.connector.id
  zcomponent_id         = zpa_app_connector_group.canada_connector_group.id
  depends_on            = [ data.zpa_enrollment_cert.connector, zpa_app_connector_group.us_connector_group]
}

# Create an App Connector Group
resource "zpa_app_connector_group" "usa_connector_group" {
  name                          = "USA Connector Group"
  description                   = "USA Connector Group"
  enabled                       = true
  city_country                  = "San Jose, CA"
  country_code                  = "CA"
  latitude                      = "43.6532"
  longitude                     = "79.3832"
  location                      = "San Jose, CA, USA"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  dns_query_type                = "IPV4"
}
```

## Service Edge Provisioning KeyExample Usage

```hcl
# Create Provisioning Key for Service Edge Group
resource "zpa_provisioning_key" "test_provisioning_key" {
  name                  = "test-provisioning-key"
  association_type      = "SERVICE_EDGE_GRP"
  max_usage             = "10"
  enrollment_cert_id    = data.zpa_enrollment_cert.service_edge.id
  zcomponent_id         = zpa_service_edge_group.service_edge_group_nyc.id
}

# Retrieve the Service Edge Enrollment Certificate
data "zpa_enrollment_cert" "service_edge" {
    name = "Service Edge"
}

# Create a Service Edge Group
resource "zpa_service_edge_group" "service_edge_group_nyc" {
  name                  = "Service Edge Group New York"
  description           = "Service Edge Group New York"
  upgrade_day           = "SUNDAY"
  upgrade_time_in_secs  = "66600"
  latitude              = "40.7128"
  longitude             = "-73.935242"
  location              = "New York, NY, USA"
  version_profile_id    = "0"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the provisioning key.
* `max_usage` - (Required) The maximum number of instances where this provisioning key can be used for enrolling an App Connector or Service Edge.
* `enrollment_cert_id` - (Required) ID of the enrollment certificate that can be used for this provisioning key. `ID` of the existing enrollment certificate that has the private key
* `zcomponent_id` - (Required) ID of the existing App Connector or Service Edge Group.
* `association_type` (Required) Specifies the provisioning key type for App Connectors or ZPA Private Service Edges. The supported values are `CONNECTOR_GRP` and `SERVICE_EDGE_GRP`
* `microtenant_id` (Optional) The ID of the microtenant the resource is to be associated with.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Provisioning key can be imported by using `<PROVISIONING KEY ID>` or `<PROVISIONING KEY NAME>` as the import ID.

For example:

```shell
terraform import zpa_provisioning_key.example <provisioning_key_id>
```

or

```shell
terraform import zpa_provisioning_key.example <provisioning_key_name>
```
