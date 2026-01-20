---
page_title: "zpa_risk_score_values Resource - terraform-provider-zpa"
subcategory: "Policy Set Controller"
description: |-
  Get information about risk scores for the specified customer.
---

# zpa_risk_score_values (Data Source)

Use the **zpa_risk_score_values** data source to get information about risk score values for the specified customer in the Zscaler Private Access cloud. This data source can be optionally used when defining the following policy types:
    - ``zpa_policy_access_rule``

## Example Usage

```terraform
data "zpa_risk_score_values" "this" {
}
```

## Schema

### Optional

- `exclude_unknown` - (String) Exclude unknown risk values

### Read-Only

The following values are returned:

- `CRITICAL`
- `HIGH`
- `MEDIUM`
- `LOW`
- `UNKNOWN`
