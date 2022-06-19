---
subcategory: "Browser Access Certificate"
layout: "zscaler"
page_title: "ZPA: ba_certificate"
description: |-
  Get information about ZPA Browser Access Certificate in Zscaler Private Access cloud.
---

# Data Source: zpa_ba_certificate

Use the **zpa_ba_certificate** data source to get information about a browser access certificate created in the Zscaler Private Access cloud. This data source is required when creating a browser access application segment resource.

## Example Usage

```hcl
# ZPA Browser Access Data Source
data "zpa_ba_certificate" "foo" {
  name = "example.acme.com"
}
```

```hcl
# ZPA Browser Access Data Source
data "zpa_ba_certificate" "foo" {
  id = "1234567890"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the browser access certificate to be exported.
* `id` - (Optional) The id of the browser access certificate to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `cert_chain` - (string)
* `certificate` - (string) The certificate text is in PEM format.
* `cname` - (string)
* `creation_time` - (string)
* `description` - (string)
* `issued_by` - (string)
* `issued_to` - (string)
* `modified_time` - (string)
* `modifiedby` - (string)
* `san` - (string)
* `serial_no` - (string)
* `status` - (string)
* `valid_from_in_epochsec` - (string)
* `valid_to_in_epochsec` - (string)

:warning: Notice that certificate and public_keys are omitted from the output.
