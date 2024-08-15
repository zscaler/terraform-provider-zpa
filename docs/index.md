---
layout: "zscaler"
page_title: "Provider: Zscaler Private Access (ZPA)"
description: |-
   The Zscaler Private Access provider is used to interact with ZPA API, to onboard new application segments, segment groups, server groups, application servers and create zero trust access policies. To use this  provider, you must create ZPA API credentials.
---

# Zscaler Private Access Provider (ZPA)

The Zscaler Private Access (ZPA) provider is used to interact with [ZPA](https://www.zscaler.com/products/zscaler-private-access) platform, to onboard new application segments, segment groups, server groups, and create zero trust access policies. To use this  provider, you must create ZPA API credentials. For details on API credentials, please visit the official product [help portal](https://help.zscaler.com/zpa/about-api-keys)

Use the navigation on the left to read about the available resources.

## Support Disclaimer

-> **Disclaimer:** Please refer to our [General Support Statement](guides/support.md) before proceeding with the use of this provider. You can also refer to our [troubleshooting guide](guides/troubleshooting.md) for guidance on typical problems.

## Zenith Community - ZPA Terraform Provider Introduction

[![ZPA Terraform provider Video Series Ep1](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_terraform_provider_introduction.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEpCAI/video-zpa-terraform-provider-video-series-ep1)

## Example Usage ZPA Production Cloud

For customers running this provider in their production tenant, the variable `ZPA_CLOUD` is optional. If provided, it must be followed by the value `PRODUCTION`.

```terraform
# Configure ZPA provider source and version
terraform {
  required_providers {
    zpa = {
      source = "zscaler/zpa"
      version = "~> 3.0.0"
    }
  }
}

provider "zpa" {
  zpa_client_id         = "xxxxxxxxxxxxxxxx"
  zpa_client_secret     = "xxxxxxxxxxxxxxxx"
  zpa_customer_id       = "xxxxxxxxxxxxxxxx"
}

resouce "zpa_application_segment" "this" {
  # ...
}
```

## Example Usage ZPA Beta, GOV, GOVUS, Preview, and Dev Cloud

For customers who want to use this provider with ZPA Beta, Gov, Preview, and Dev Cloud, the following variable credentials `zpa_cloud` followed by the value `BETA`, `ZPATWO`, `GOV`, `GOVUS`, or `PREVIEW` values or via environment variable `ZPA_CLOUD=BETA`, `ZPA_CLOUD=ZPATWO`, `ZPA_CLOUD=GOV`, `ZPA_CLOUD=GOVUS`, `ZPA_CLOUD=PREVIEW`, `ZPA_CLOUD=DEV`are required.

```terraform
# Configure ZPA provider source and version
terraform {
  required_providers {
    zpa = {
      source = "zscaler/zpa"
      version = "~> 3.0.0"
    }
  }
}

provider "zpa" {
  zpa_client_id         = "xxxxxxxxxxxxxxxx"
  zpa_client_secret     = "xxxxxxxxxxxxxxxx"
  zpa_customer_id       = "xxxxxxxxxxxxxxxx"
  zpa_cloud             = "BETA" // Use `BETA`, `GOV`, `GOVUS`, `PREVIEW` or `DEV`
}

resouce "zpa_application_segment" "app_segment" {
  # ...
}
```

## Terraform / Zscaler Private Access Interaction

### Parallelism

Terraform uses goroutines to speed up deployment, but the number of parallel
operations it launches may exceed [what is recommended](https://help.zscaler.com/zpa/about-rate-limiting).
When configuring ZPA Policies we recommend to limit the number of concurrent API calls to **ONE**. This limit ensures that there is no performance impact during the provisioning of large Terraform configurations involving access policy creation.

This recommendation applies to the following resources:

- ``zpa_policy_access_rule``
- ``zpa_policy_inspection_rule``
- ``zpa_policy_timeout_rule``
- ``zpa_policy_forwarding_rule``
- ``zpa_policy_isolation_rule``

In order to accomplish this, we recommend setting the [parallelism](https://www.terraform.io/cli/commands/apply#parallelism-n) value at this limit to prevent performance impacts.

## Authentication

The ZPA provider offers various means of providing credentials for authentication. The following methods are supported:

* Directly in the provider block
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

~> **NOTE** `ZPA_CLOUD` environment variable is an optional parameter when running this provider in production; however, this parameter is required to provision resources in the ZPA Beta Cloud, Gov Cloud, Gov US Cloud, or Preview Cloud.

~> **NOTE** `ZPA_MICROTENANT_ID` environment variable is an optional parameter when provisioning resources within a ZPA microtenant

```terraform
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
$env:ZPA_CLIENT_ID='xxxxxxxxxxxxxxxx'
$env:ZPA_CLIENT_SECRET='xxxxxxxxxxxxxxxx'
$env:ZPA_CUSTOMER_ID='xxxxxxxxxxxxxxxx'
terraform plan
```

### Configuration file

You can use a configuration file to specify your credentials. The
file location must be `$HOME/.zpa/credentials.json` on Linux and OS X, or
`"%USERPROFILE%\.zpa/credentials.json"` for Windows users.
If we fail to detect credentials inline, or in the environment variable, Terraform will check
this location.

Usage:

```terraform
provider "zpa" {}
```

credentials.json file:

```json
{
  "zpa_client_id":"zpa_client_id",
  "zpa_client_secret": "zpa_client_secret",
  "zpa_customer_id": "zpa_customer_id"
}
```

## Argument Reference

The following arguments are supported:

### Required

* ``zpa_client_id`` - (Required) ZPA client ID, is equivalent to a username.
* ``zpa_client_secret`` - (Required) ZPA client secret, is equivalent to a secret password.
* ``zpa_customer_id`` - (Required) ZPA customer ID, is equivalent to your ZPA tenant identification.
* ``zpa_cloud`` - (Required) ZPA Cloud name `BETA`, `GOV`, `GOVUS`, `PREVIEW` or `ZPATWO` clouds.

### Optional

* `zpa_cloud` - (Optional) ZPA Cloud name `PRODUCTION`. Optional when running in the ZPA production cloud.

### Zscaler Private Access Microtenant

A Microtenant is a delegated administrator responsibility that is assigned to an admin by an admin with Microtenant administrator privileges. Microtenants are defined by an authentication domain and assigned to admins based on country, department, and company for role-based administration control. A Microtenant is created within a tenant and is used when departments or subsidiaries within an organization want to manage their configurations independently.[Read More](https://help.zscaler.com/zpa/about-microtenants)

To manage a microtenant using the ZPA Terraform provider, the administrator for the parent or default tenant, must first provision the microtenant using the resource `zpa_microtenant_controller`. The resource will output the administrator credentials for the new microtenant, which can then be provided to the microtenant administrator.

The microtenant administrator can then create his own microtenant API credentials required to authenticate via API to the ZPA platform. From that point, the administrator can then individually manage his own resources in an isolated manner.
When authenticating to microtenant via API using the ZPA Terraform provider, the administrator must provide the following environment variable credentials: `ZPA_CLIENT_ID`, `ZPA_CLIENT_SECRET`, `ZPA_CUSTOMER_ID`, `ZPA_CLOUD`, `ZPA_MICROTENANT_ID`

~> **NOTE 1** Only environment variables are currently supported when authenticating to a Microtenant.

~> **NOTE 2** The environment variable `ZPA_MICROTENANT_ID` is mandatory when provisioning/managing resources exclusively within a Microtenant.

~> **NOTE 3** `ZPA_CLOUD` environment variable is an optional parameter when running this provider in production; however, this parameter is required to provision resources in the `BETA`, `GOV`, `GOVUS`, `PREVIEW` or `ZPATWO` clouds.

## Support Disclaimer

-> **Disclaimer:** Please refer to our [General Support Statement](guides/support.md) before proceeding with the use of this provider. You can also refer to our [troubleshooting guide](guides/troubleshooting.md) for guidance on typical problems.
