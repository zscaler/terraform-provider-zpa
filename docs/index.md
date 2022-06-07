---
layout: "zscaler"
page_title: "Provider: Zscaler Private Access (ZPA)"
description: |-
   The Zscaler Private Access provider is used to interact with ZPA API, to onboard new application segments, segment groups, server groups, application servers and create zero trust access policies. To use this  provider, you must create ZPA API credentials.
---

# Zscaler Private Access Provider (ZPA)

The Zscaler Private Access (ZPA) provider is used to interact with [ZPA](https://www.zscaler.com/products/zscaler-private-access) platform, to onboard new application segments, segment groups, server groups, and create zero trust access policies. To use this  provider, you must create ZPA API credentials. For details on API credentials, please visit the official product [help portal](https://help.zscaler.com/zpa/about-api-keys)

Use the navigation on the left to read about the available resources.

## Example Usage ZPA Production Cloud

For customers running this provider in their production tenant, the variable `ZPA_CLOUD` is optional. If provided, it must be followed by the value `PRODUCTION`.

```hcl
# Configure ZPA provider source and version
terraform {
  required_providers {
    zpa = {
      source = "zscaler/zpa"
      version = "2.1.6"
    }
  }
}

provider "zpa" {
  zpa_client_id         = "xxxxxxxxxxxxxxxx"
  zpa_client_secret     = "xxxxxxxxxxxxxxxx"
  zpa_customer_id       = "xxxxxxxxxxxxxxxx"
}

resouce "zpa_application_segment" "app_segment" {
  # ...
}
```

## Example Usage ZPA Beta Cloud

For customers who want to use this provider with ZPA Beta Cloud, they must provide the following variable credentials `zpa_cloud` followed by the value `BETA` or via environment variable `ZPA_CLOUD=BETA`.

```hcl
# Configure ZPA provider source and version
terraform {
  required_providers {
    zpa = {
      source = "zscaler/zpa"
      version = "2.1.6"
    }
  }
}

provider "zpa" {
  zpa_client_id         = "xxxxxxxxxxxxxxxx"
  zpa_client_secret     = "xxxxxxxxxxxxxxxx"
  zpa_customer_id       = "xxxxxxxxxxxxxxxx"
  zpa_cloud             = "BETA"
}

resouce "zpa_application_segment" "app_segment" {
  # ...
}
```

## Authentication

The ZPA provider offers various means of providing credentials for authentication. The following methods are supported:

* Static credentials directly in the provider's block
* Environment variables
* From the JSON config file

### Static credentials

!> **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

Static credentials can be provided by specifying the `zpa_client_id`, `zpa_client_secret` and `zpa_customer_id` arguments in-line in the ZPA provider block:

**Usage:**

``` hcl
provider "zpa" {
  zpa_client_id         = "xxxxxxxxxxxxxxxx"
  zpa_client_secret     = "xxxxxxxxxxxxxxxx"
  zpa_customer_id       = "xxxxxxxxxxxxxxxx"
}
```

### Environment variables

You can provide credentials via the `ZPA_CLIENT_ID`, `ZPA_CLIENT_SECRET`, `ZPA_CUSTOMER_ID`, `ZPA_CLOUD` environment variables, representing your ZPA API key credentials and customer ID, of your ZPA account, respectively.

~> **NOTE** `ZPA_CLOUD` environment variable is option when running this provider in production, but required if running in the ZPA Beta Cloud.

```hcl
provider "zpa" {}
```

**macOS and Linux Usage:**

```sh
export ZPA_CLIENT_ID      = "xxxxxxxxxxxxxxxx"
export ZPA_CLIENT_SECRET  = "xxxxxxxxxxxxxxxx"
export ZPA_CUSTOMER_ID    = "xxxxxxxxxxxxxxxx"
terraform plan
```

**Windows Powershell:**

```powershell
env:ZPA_CLIENT_ID      = 'xxxxxxxxxxxxxxxx'
env:ZPA_CLIENT_SECRET  = 'xxxxxxxxxxxxxxxx'
env:ZPA_CUSTOMER_ID    = 'xxxxxxxxxxxxxxxx'
terraform plan
```

### Configuration file

You can use a configuration file to specify your credentials. The
file location is `$HOME/.zpa/credentials.json` on Linux and OS X, or
`"%USERPROFILE%\.zpa/credentials.json"` for Windows users.
If we fail to detect credentials inline, or in the environment, Terraform will check
this location.

Usage:

```hcl
provider "zpa" {}
```

credentials.json file:

```json
{
  "zpa_client_id": "XXXXXXXXXX",
  "zpa_client_secret": "XXXXXXXXXX",
  "zpa_customer_id": "XXXXXXXXXX"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `zpa_client_id` - (Required) ZPA client ID, is equivalent to a username.
* `zpa_client_secret` - (Required) ZPA client secret, is equivalent to a secret password.
* `zpa_customer_id` - (Required) ZPA customer ID, is equivalent to your ZPA tenant identification.
* `zpa_cloud` - (Required) ZPA Cloud name `BETA`. Only required when running in the ZPA beta cloud.
  * This parameter is `NOT` required when running the provider to provision the ZPA Production cloud.

### Optional

## Support

This template/solution are released under an as-is, best effort, support
policy. These scripts should be seen as community supported and Zscaler
Business Development Team will contribute our expertise as and when possible.
We do not provide technical support or help in using or troubleshooting the components
of the project through our normal support options such as Zscaler support teams,
or ASC (Authorized Support Centers) partners and backline
support options. The underlying product used (Zscaler Private Access API) by the
scripts or templates are still supported, but the support is only for the
product functionality and not for help in deploying or using the template or
script itself. Unless explicitly tagged, all projects or work posted in our
GitHub repository at (<https://github.com/zscaler>) or sites other
than our official Downloads page on (<https://help.zscaler.com/login-tickets>)
are provided under the best effort policy.
