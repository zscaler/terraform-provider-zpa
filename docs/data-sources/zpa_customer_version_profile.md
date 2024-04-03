---
page_title: "zpa_customer_version_profile Data Source - terraform-provider-zpa"
subcategory: "Customer Version Profile"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-connectors
  API documentation https://help.zscaler.com/zpa/obtaining-version-profile-details-using-api
  Get information about all customer version profile details.
---

# zpa_customer_version_profile (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-connectors)
* [API documentation](https://help.zscaler.com/zpa/obtaining-version-profile-details-using-api)

Use the **zpa_customer_version_profile** data source to get information about all customer version profiles from the Zscaler Private Access cloud. This data source can be associated with an App Connector Group within the parameter `version_profile_id` or `version_profile_name`

The customer version profile IDs are:

* `Default` = `0`
* `Previous Default` = `1`
* `New Release` = `2`

## Example Usage

```terraform
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

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the enrollment certificate to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) The id of the enrollment certificate to be exported.
* `allow_signing` - (bool)
* `cname` - (string)
* `certificate` - (string) The certificate text is in PEM format.
* `client_cert_type` - (string) Returned values are:
  * `ZAPP_CLIENT`
  * `ISOLATION_CLIENT`
  * `NONE`

* `creation_time` - (string)
* `csr` - (string)
* `description` - (string)
* `issued_by` - (string)
* `issued_to` - (string)
* `modified_time` - (string)
* `modified_by` - (string)
* `parent_cert_id` - (string)
* `parent_cert_name` - (string)
* `cert_chain` - (string)
* `serial_no` - (string)
* `valid_from_in_epoch_sec` - (string)
* `valid_to_in_epochsec` - (string)
