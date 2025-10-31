---
page_title: "zpa_application_segment_multimatch_bulk Data source - terraform-provider-zpa"
subcategory: "Application Segment"
description: |-
  Official documentation https://help.zscaler.com/zpa/using-app-segment-multimatch
  API documentation https://help.zscaler.com/zpa/configuring-application-segment-multimatch-using-api
  Bulk updates application segment Multimatch in multiple applications.
---

# zpa_application_segment_multimatch_bulk (Resource)

* [Official documentation](https://help.zscaler.com/zpa/adding-ip-ranges)
* [API documentation](https://help.zscaler.com/zpa/adding-ip-ranges)

The **zpa_application_segment_multimatch_bulk** resource to bulk updates application segment Multimatch in multiple applications.

## Example Usage

```terraform
resource "zpa_application_segment_multimatch_bulk" "this" {
  application_ids = ["72058304855164528","72058304855164536"]
  match_style = "EXCLUSIVE"
}
```

## Schema

### Required

- `application_ids` (List of Integers) The list of Application Segment IDs
- `match_style` (String) Indicates if Multimatch is enabled for the application segment. If enabled (INCLUSIVE), the request allows traffic to match multiple applications. If disabled (EXCLUSIVE), the request allows traffic to match a single application. A domain can only be INCLUSIVE or EXCLUSIVE, and any application segment can only contain inclusive or exclusive domains. Supported values: `EXCLUSIVE` and `INCLUSIVE`
