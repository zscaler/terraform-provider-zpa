---
subcategory: "Policy Set Controller"
layout: "zscaler"
page_title: "ZPA: access_policy_client_types"
description: |-
  Get information about all client types for the specified customer.
---

# Data Source: zpa_access_policy_client_types

Use the **zpa_access_policy_client_types** data source to get information about all client types for the specified customer in the Zscaler Private Access cloud. This data source can be optionally used when defining the following policy types:
    - ``zpa_policy_access_rule``
    - ``zpa_policy_timeout_rule``
    - ``zpa_policy_forwarding_rule``
    - ``zpa_policy_isolation_rule``
    - ``zpa_policy_inspection_rule``

The ``object_type`` attribute must be defined as "CLIENT_TYPE" in the policy operand condition. To learn more see the To learn more see the [Getting Details of All Client Types](https://help.zscaler.com/zpa/configuring-access-policies-using-api#getClientTypes)

-> **NOTE** By Default the ZPA provider will return all client types

-> **NOTE** When defining a ``zpa_policy_isolation_rule`` policy the ``object_type`` "CLIENT_TYPE" is mandatory and ``zpn_client_type_exporter`` is the only supported value.

## Example Usage

```hcl
data "zpa_access_policy_client_types" "this" {
}
```

## Argument Reference

The following values are returned:

* `"zpn_client_type_branch_connector" = "Branch Connector"`
* `"zpn_client_type_browser_isolation" = "Cloud Browser"`
* `"zpn_client_type_edge_connector" = "Cloud Connector"`
* `"zpn_client_type_exporter" = "Web Browser"`
* `"zpn_client_type_exporter_noauth" = "Web Browser Unauthenticated"`
* `"zpn_client_type_ip_anchoring" = "ZIA Service Edge"`
* `"zpn_client_type_machine_tunnel" = "Machine Tunnel"`
* `"zpn_client_type_slogger" = "ZPA LSS"`
* `"zpn_client_type_zapp" = "Client Connector"`

To learn more see the [Getting Details of All Client Types](https://help.zscaler.com/zpa/configuring-access-policies-using-api#getClientTypes)
