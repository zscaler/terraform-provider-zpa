---
subcategory: "Customer Version Profile"
layout: "zscaler"
page_title: "ZPA: customer_version_profile"
description: |-
  Get information about all configured enrollment certificate details.
---

# zpa_customer_version_profile

Use the **zpa_customer_version_profile** data source to get information about all customer version profiles from the Zscaler Private Access cloud. This data source can be associated with an App Connector Group within the parameter `version_profile_id` or `version_profile_name`

The customer version profile IDs are:

* `Default` = `0`
* `Previous Default` = `1`
* `New Release` = `2`

## Example Usage

```hcl
# Retrieve "Default" customer version profile
data "zpa_customer_version_profile" "default" {
    name = "Default"
}

# Retrieve "Previous Default" customer version profile
data "zpa_customer_version_profile" "previous_default"{
    name = "Previous Default"
}

# Retrieve "New Release" customer version profile
data "zpa_customer_version_profile" "new_release"{
    name = "New Release"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the enrollment certificate to be exported.
* `id` - (Optional) The id of the enrollment certificate to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `allow_signing` - (Boolean)
* `cname` - (String)
* `certificate` - (String) The certificate text is in PEM format.
* `client_cert_type` - (String) Returned values are:
  * `ZAPP_CLIENT`
  * `ISOLATION_CLIENT`
  * `NONE`

* `creation_time` - (String)
* `csr` - (String)
* `description` - (String)
* `issued_by` - (String)
* `issued_to` - (String)
* `modified_time` - (String)
* `modifiedby` - (String)
* `parent_cert_id` - (String)
* `parent_cert_name` - (String)
* `cert_chain` - (String)
* `serial_no` - (String)
* `valid_from_in_epoch_sec` - (String)
* `valid_to_in_epochsec` - (String)
