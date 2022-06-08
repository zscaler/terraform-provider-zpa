---
subcategory: "Policy Set Controller"
layout: "zscaler"
page_title: "ZPA: policy_type"
description: |-
  Get information about Policy Set ID in Zscaler Private Access cloud.
---

# Data Source: zpa_policy_timeout

Use the **zpa_policy_type** data source to get information about an a ``policy_set_id`` and ``policy_type``. This data source is required when creating:

1. Access policy Rules
2. Access policy timeout rules
3. Access policy forwarding rules
4. Access policy inspection rules
5. Access policy isolation rules

~> **NOTE** The parameters ``policy_set_id`` is required in all circumstances and is exported when checking for the policy_type parameter. The policy_type value is used for differentiating the policy types, in the request endpoint. The supported values are:

* ``ACCESS_POLICY/GLOBAL_POLICY``
* ``TIMEOUT_POLICY/REAUTH_POLICY``
* ``BYPASS_POLICY/CLIENT_FORWARDING_POLICY``
* ``INSPECTION_POLICY``
* ``ISOLATION_POLICY``
* ``SIEM_POLICY``

## Example Usage

```hcl
# Get information for "ACCESS_POLICY" ID
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

output "zpa_policy_type_access_policy" {
    value = data.zpa_policy_type.access_policy.id
}
```

```hcl
# Get information for "GLOBAL_POLICY" ID
data "zpa_policy_type" "global_policy" {
    policy_type = "GLOBAL_POLICY"
}

output "zpa_policy_type_access_policy" {
    value = data.zpa_policy_type.global_policy.id
}
```

```hcl
# Get information for "TIMEOUT_POLICY" ID
data "zpa_policy_type" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}

output "zpa_policy_type_timeout_policy" {
    value = data.zpa_policy_type.timeout_policy.id
}
```

```hcl
# Get information for "REAUTH_POLICY" ID
data "zpa_policy_type" "reauth_policy" {
    policy_type = "REAUTH_POLICY"
}

output "zpa_policy_type_reauth_policy" {
    value = data.zpa_policy_type.reauth_policy.id
}
```

```hcl
# Get information for "CLIENT_FORWARDING_POLICY" ID
data "zpa_policy_type" "client_forwarding_policy" {
    policy_type = "CLIENT_FORWARDING_POLICY"
}

output "zpa_policy_type_client_forwarding_policy" {
    value = data.zpa_policy_type.client_forwarding_policy.id
}
```

```hcl
# Get information for "INSPECTION_POLICY" ID
data "zpa_policy_type" "inspection_policy" {
    policy_type = "INSPECTION_POLICY"
}

output "zpa_policy_type_inspection_policy" {
    value = data.zpa_policy_type.inspection_policy.id
}
```

```hcl
# Get information for "ISOLATION_POLICY" ID
data "zpa_policy_type" "isolation_policy" {
    policy_type = "ISOLATION_POLICY"
}

output "zpa_policy_type_isolation_policy" {
    value = data.zpa_policy_type.isolation_policy.id
}
```

```hcl
# Get information for "SIEM_POLICY" ID
data "zpa_policy_type" "siem_policy" {
    policy_type = "SIEM_POLICY"
}

output "zpa_policy_type_siem_policy" {
    value = data.zpa_policy_type.siem_policy
}
```

## Argument Reference

The following arguments are supported:

* `policy_type` - (Optional) The value for differentiating the policy types.
* `policy_set_id` - (Required) The ID of the global policy set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* ``creation_time`` (string)
* ``description`` (string)
* ``enabled``  (bool)
* ``id`` (string)
* ``modified_time``  (string)
* ``modified_by``  (string)
* ``name``  (string)
* ``sorted`` (bool)
