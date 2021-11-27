---
layout: "ZPA"
page_title: "Provider: Zscaler Private Access"
description: |-
   The Zscaler Private Access provider is used to interact with ZPA API, to onboard new application segments, segment groups, server groups, application servers and create zero trust access policies. To use this  provider, you must create ZPA API credentials.

---

⚠️ **Attention:** This provider is not affiliated with, nor supported by Zscaler in any way.

# Zscaler Private Access (ZPA) Provider

The Zscaler Private Access (ZPA) provider is used to interact with [ZPA](https://www.zscaler.com/products/zscaler-private-access) platform, to onboard new application segments, segment groups, server groups, and create zero trust access policies. To use this  provider, you must create ZPA API credentials. For details on API credentials, please visit the official product [help portal](https://help.zscaler.com/zpa/about-api-keys)


Use the navigation on the left to read about the available resources.

## Authentication

The ZPA provider offers various means of providing credentials for authentication. The following methods are supported:

* Static credentials
* Environment variables

### Static credentials

⚠️ **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

Static credentials can be provided by specifying the `zpa_client_id`, `zpa_client_secret` and `zpa_customer_id` arguments in-line in the ZPA provider block:

**Usage:**

```hcl
provider "zpa" {
  zpa_client_id         = "xxxxxxxxxxxxxxxx"
  zpa_client_secret     = "xxxxxxxxxxxxxxxx"
  zpa_customer_id       = "xxxxxxxxxxxxxxxx"
}
```

### Environment variables

You can provide credentials via the `ZPA_CLIENT_ID`, `ZPA_CLIENT_SECRET`, `ZPA_CUSTOMER_ID` environment variables, representing your ZPA API key credentials and customer ID, of your ZPA account, respectively.

```hcl
provider "zpa" {}
```

**macOS and Linux Usage:**

```sh
$ export ZPA_CLIENT_ID      = "xxxxxxxxxxxxxxxx"
$ export ZPA_CLIENT_SECRET  = "xxxxxxxxxxxxxxxx"
$ export ZPA_CUSTOMER_ID    = "xxxxxxxxxxxxxxxx"
$ terraform plan
```

**Windows Powershell:**

```powershell
$env:ZPA_CLIENT_ID      = 'xxxxxxxxxxxxxxxx'
$env:ZPA_CLIENT_SECRET  = 'xxxxxxxxxxxxxxxx'
$env:ZPA_CUSTOMER_ID    = 'xxxxxxxxxxxxxxxx'
```

## Argument Reference

The following arguments are supported:

### Required

* `zpa_client_id` - (Required) ZPA client ID, is equivalent to a username.
* `zpa_client_secret` - (Required) ZPA client secret, is equivalent to a secret password.
* `zpa_customer_id` - (Required) ZPA customer ID, is equivalent to your ZPA tenant identification.
