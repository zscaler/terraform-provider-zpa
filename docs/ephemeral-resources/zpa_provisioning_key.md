---
page_title: "zpa_provisioning_key Ephemeral Resource - terraform-provider-zpa"
subcategory: "Provisioning Key"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-connector-provisioning-keys
  API documentation https://help.zscaler.com/zpa/configuring-provisioning-keys-using-api
  Creates and manages ZPA Provisioning Key for Service Edge and/or App Connector Groups.
---

# zpa_provisioning_key (Ephemeral Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-connector-provisioning-keys)
* [API documentation](https://help.zscaler.com/zpa/configuring-provisioning-keys-using-api)

The **zpa_provisioning_key** ephemeral resource retrieves the provisioning key value without writing it to Terraform state. Combine it with the managed `zpa_provisioning_key` resource when you need the raw key during apply.

> ℹ️ **Usage pattern:** Create or reference a `zpa_provisioning_key` resource, then add an `ephemeral "zpa_provisioning_key"` block pointing at its `id` and `association_type`. During `terraform apply`, the ephemeral block fetches the key and makes it available only for the duration of the run.
>
> ⚠️ **Consume immediately:** Because Terraform discards ephemeral values after the run, you must handle the key inside the same apply execution (for example, by writing it to a secure file, pushing it to a secrets manager, etc.). Root-module outputs cannot expose ephemerals.

## Zenith Community - ZPA Provisioning Keys

[![ZPA Terraform provider Video Series Ep3 - Provisioning Keys](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_provisioning_key.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEnCAI/video-zpa-terraform-provider-video-series-ep3-provisioning-keys)

## App Connector Group Provisioning Key Example Usage

```terraform
# Retrieve the Connector Enrollment Certificate
data "zpa_enrollment_cert" "connector" {
    name = "Connector"
}

# Create an App Connector Group
resource "zpa_app_connector_group" "this" {
  name                          = "USA Connector Group"
  description                   = "USA Connector Group"
  enabled                       = true
  city_country                  = "San Jose, US"
  country_code                  = "US"
  latitude                      = "37.33874"
  longitude                     = "-121.8852525"
  location                      = "San Jose, CA, USA"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = "0"
  dns_query_type                = "IPV4"
}

# Create Provisioning Key for App Connector Group
resource "zpa_provisioning_key" "this" {
  name                  = "test_provisioning_key"
  association_type      = "CONNECTOR_GRP"
  max_usage             = "10"
  enrollment_cert_id    = data.zpa_enrollment_cert.connector.id
  zcomponent_id         = zpa_app_connector_group.this.id
  depends_on            = [ data.zpa_enrollment_cert.connector, zpa_app_connector_group.this]
}

ephemeral "zpa_provisioning_key" "this" {
  id               = zpa_provisioning_key.this.id
  association_type = zpa_provisioning_key.this.association_type
}

# Persist the ephemeral value outside Terraform state (optional example)
resource "null_resource" "write_provisioning_key" {
  triggers = {
    apply_time = time_static.apply_time.rfc3339
  }

  provisioner "local-exec" {
    command = <<-EOC
      printf %s "${ephemeral.zpa_provisioning_key.this.provisioning_key}" > "${path.module}/provisioning_key.txt"
    EOC
  }
}

resource "time_static" "apply_time" {}

# Although this example writes the key to a local file, the value is never persisted to Terraform state.
# For production usage, prefer pushing the key into a secure secrets manager such as HashiCorp Vault
# or AWS Secrets Manager during the same apply run.
```

## Service Edge Provisioning KeyExample Usage

```terraform
# Retrieve the Service Edge Enrollment Certificate
data "zpa_enrollment_cert" "service_edge" {
    name = "Service Edge"
}

# Create a Service Edge Group
resource "zpa_service_edge_group" "this" {
  name                  = "Service Edge Group New York"
  description           = "Service Edge Group New York"
  upgrade_day           = "SUNDAY"
  upgrade_time_in_secs  = "66600"
  latitude              = "40.7128"
  longitude             = "-73.935242"
  location              = "New York, NY, USA"
  version_profile_id    = "0"
}

# Create Provisioning Key for Service Edge Group
resource "zpa_provisioning_key" "this" {
  name                  = "test-provisioning-key"
  association_type      = "SERVICE_EDGE_GRP"
  max_usage             = "10"
  enrollment_cert_id    = data.zpa_enrollment_cert.service_edge.id
  zcomponent_id         = zpa_service_edge_group.this.id
}

ephemeral "zpa_provisioning_key" "this" {
  id               = zpa_provisioning_key.this.id
  association_type = zpa_provisioning_key.this.association_type
}

resource "null_resource" "write_service_edge_key" {
  triggers = {
    apply_time = time_static.apply_time.rfc3339
  }

  provisioner "local-exec" {
    command = <<-EOC
      printf %s "${ephemeral.zpa_provisioning_key.this.provisioning_key}" > "${path.module}/service_edge_provisioning_key.txt"
    EOC
  }
}

resource "time_static" "apply_time" {}
```

## Schema

### Required

The following arguments are supported:


* `association_type` (String) Specifies the provisioning key type for App Connectors or ZPA Private Service Edges. The supported values are `CONNECTOR_GRP` and `SERVICE_EDGE_GRP`.

### Optional

* `microtenant_id` (String) The ID of the microtenant the resource is associated with.
* `provisioning_key` (String, Sensitive) The provisioning key returned by the API. Available only during apply; never persisted to Terraform state. Consume it immediately (for example, using `local-exec`, Vault, or another secure sink).

⚠️ **WARNING:** The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID`, which is the recommended method.

