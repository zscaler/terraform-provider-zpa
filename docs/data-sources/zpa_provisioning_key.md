---
page_title: "zpa_provisioning_key Resource - terraform-provider-zpa"
subcategory: "Provisioning Key"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-connector-provisioning-keys
  API documentation https://help.zscaler.com/zpa/configuring-provisioning-keys-using-api
  Get information about Provisioning Key in Zscaler Private Access cloud.
---

# Data Source: zpa_provisioning_key

* [Official documentation](https://help.zscaler.com/zpa/about-connector-provisioning-keys)
* [API documentation](https://help.zscaler.com/zpa/configuring-provisioning-keys-using-api)

Use the **zpa_provisioning_key** data source to get information about a provisioning key in the Zscaler Private Access portal or via API. This data source can be referenced in the following ZPA resources:

* App Connector Groups
* Service Edge Groups

-> **NOTE** The ``association_type`` parameter is required in order to distinguish between ``CONNECTOR_GRP`` and ``SERVICE_EDGE_GRP``

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Zenith Community - ZPA Provisioning Keys

[![ZPA Terraform provider Video Series Ep3 - Provisioning Keys](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_provisioning_key.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEnCAI/video-zpa-terraform-provider-video-series-ep3-provisioning-keys)

## Example Usage

### Basic Usage

```terraform
# ZPA Provisioning Key for "CONNECTOR_GRP"
data "zpa_provisioning_key" "connector_key" {
  name             = "Connector_Provisioning_Key"
  association_type = "CONNECTOR_GRP"
}
```

```terraform
# ZPA Provisioning Key for "SERVICE_EDGE_GRP"
data "zpa_provisioning_key" "service_edge_key" {
  name             = "ServiceEdge_Provisioning_Key"
  association_type = "SERVICE_EDGE_GRP"
}
```

### Accessing the Provisioning Key Value

The provisioning key value is marked as sensitive and can be accessed using outputs or resource references:

```terraform
# Retrieve existing provisioning key
data "zpa_provisioning_key" "existing" {
  name             = "Production_Connector_Key"
  association_type = "CONNECTOR_GRP"
}

# Output the provisioning key (marked as sensitive)
output "provisioning_key_value" {
  description = "Use this key to onboard App Connectors"
  value       = data.zpa_provisioning_key.existing.provisioning_key
  sensitive   = true
}

# Use in automation
resource "null_resource" "deploy_connector" {
  provisioner "local-exec" {
    command = "deploy-connector.sh ${data.zpa_provisioning_key.existing.provisioning_key}"
  }
}
```

To retrieve the key value:
```bash
# View the provisioning key
terraform output provisioning_key_value

# Or get it programmatically
terraform output -json provisioning_key_value | jq -r .
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) Name of the provisioning key.
* `association_type` (Required) Specifies the provisioning key type for App Connectors or ZPA Private Service Edges. The supported values are `CONNECTOR_GRP` and `SERVICE_EDGE_GRP`
* `id` - (Optional) The ID of the provisioning key to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The unique identifier of the provisioning key
* `creation_time` - (String) Timestamp when the key was created
* `enabled` - (Boolean) Whether the provisioning key is enabled
* `expiration_in_epoch_sec` - (String) Expiration time in epoch seconds
* `ip_acl` - (Set of String) List of IP addresses or CIDR ranges allowed to use this key
* `max_usage` - (String) Maximum number of times this key can be used
* `modified_by` - (String) ID of the user who last modified the key
* `modified_time` - (String) Timestamp when the key was last modified
* `provisioning_key` - (String, **Sensitive**) **The actual provisioning key value**. This is the key needed to onboard App Connector or Service Edge devices. Marked as sensitive to prevent exposure in logs and console output. Access it using Terraform outputs or resource references as shown in the examples above.
* `enrollment_cert_id` - (String) ID of the enrollment certificate
* `enrollment_cert_name` - (String) Name of the enrollment certificate
* `ui_config` - (String) UI configuration for the key
* `usage_count` - (String) Number of times the key has been used
* `zcomponent_id` - (String) ID of the associated App Connector or Service Edge Group
* `zcomponent_name` - (String) Name of the associated component
* `app_connector_group_id` - (String) ID of the App Connector Group (if applicable)
* `app_connector_group_name` - (String) Name of the App Connector Group (if applicable)
* `microtenant_id` - (String) The ID of the microtenant the resource is associated with
* `microtenant_name` - (String) The name of the microtenant the resource is associated with

⚠️ **Security Note:** The `provisioning_key` attribute is stored in the Terraform state file. Ensure your state backend is properly secured with encryption at rest and access controls. See the [resource documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_provisioning_key) for detailed security guidance.
