---
page_title: "Troubleshooting Guide"
---

# How to troubleshoot your problem

If you have problems with code that uses ZPA Terraform provider, follow these steps to solve them:

* Check symptoms and solutions in the [Typical problems](#typical-problems) section below.
* Upgrade provider to the latest version. The bug might have already been fixed.
* In case of authentication problems, see the [Authentication Issues](#authentication-issues) below.
* Collect debug information using following command:

```sh
TF_LOG=DEBUG ZSCALER_SDK_VERBOSE=true ZSCALER_SDK_LOG=true terraform apply -no-color 2>&1 |tee tf-debug.log
```

* Open a [new GitHub issue](https://github.com/zscaler/terraform-provider-zpa/issues/new/choose) providing all information described in the issue template - debug logs, your Terraform code, Terraform & plugin versions, etc.

## Typical problems

### Authentication Issues

### │ Error: Invalid provider configuration and Error: failed configuring the provider

The most common problem with invalid provider is when the ZPA API credentials are not properly set via one of the supported methods. Please make sure to read the documentation for the supported authentication methods [Authentication Methods](https://registry.terraform.io/providers/zscaler/zpa/latest/docs)

```sh
│ Provider "zscaler/zpa" requires explicit configuration. Add a provider block to the root module and configure the
│ provider's required arguments as described in the provider documentation.
```

```sh
│ Error: failed configuring the provider
│
│   with provider["zscaler/zpa"],
│   on <empty> line 0:
│   (source code not available)
│
│ error:Could not open credentials file, needs to contain one json object with keys: zpa_client_id, zpa_client_secret,
│ zpa_customer_id, and zpa_cloud. open /Users/username/.zpa/credentials.json: no such file or directory
```

## Multiple Provider Configurations

The most common reason for technical difficulties might be related to missing `alias` attribute in `provider "zpa" {}` blocks or `provider` attribute in `resource "zpa_..." {}` blocks, when using multiple provider configurations. Please make sure to read [`alias`: Multiple Provider Configurations](https://www.terraform.io/docs/language/providers/configuration.html#alias-multiple-provider-configurations) documentation article.

## Error while installing: registry does not have a provider

```sh
Error while installing hashicorp/zpa: provider registry
registry.terraform.io does not have a provider named
registry.terraform.io/hashicorp/zpa
```

If you notice below error, it might be due to the fact that [required_providers](https://www.terraform.io/docs/language/providers/requirements.html#requiring-providers) block is not defined in *every module*, that uses ZPA Terraform Provider. Create `versions.tf` file with the following contents:

```hcl
# versions.tf
terraform {
  required_providers {
    zpa = {
      source  = "zscaler/zpa"
      version = "2.82.4"
    }
  }
}
```

... and copy the file in every module in your codebase. Our recommendation is to skip the `version` field for `versions.tf` file on module level, and keep it only on the environment level.

```
├── environments
│   ├── sandbox
│   │   ├── README.md
│   │   ├── main.tf
│   │   └── versions.tf
│   └── production
│       ├── README.md
│       ├── main.tf
│       └── versions.tf
└── modules
    ├── first-module
    │   ├── ...
    │   └── versions.tf
    └── second-module
        ├── ...
        └── versions.tf
```

## Error: Failed to install provider

Running the `terraform init` command, you may see `Failed to install provider` error if you didn't check-in [`.terraform.lock.hcl`](https://www.terraform.io/language/files/dependency-lock#lock-file-location) to the source code version control:

```sh
Error: Failed to install provider

Error while installing zscaler/zpa: v2.82.0: checksum list has no SHA-256 hash for "https://github.com/zscaler/terraform-provider-zpa/releases/download/v2.82.0/terraform-provider-zpa_2.82.0_darwin_amd64.zip"
```

You can fix it by following three simple steps:

* Replace `zscaler.com/zpa/zpa` with `zscaler/zpa` in all your `.tf` files with the `python3 -c "$(curl -Ls https://github.com/zscaler/terraform-provider-zpa/scripts/upgrade-namespace.py)"` command.
* Run the `terraform state replace-provider zscaler.com/zpa/zpa zscaler/zpa` command and approve the changes. See [Terraform CLI](https://www.terraform.io/cli/commands/state/replace-provider) docs for more information.
* Run `terraform init` to verify everything working.

The terraform apply command should work as expected now.

## Error: Failed to query available provider packages

See the same steps as in [Error: Failed to install provider](#error-failed-to-install-provider).

### Error: Provider registry.terraform.io/zscaler/zpa v... does not have a package available for your current platform, windows_386

This kind of errors happens when the 32-bit version of ZPA Terraform provider is used, usually on Microsoft Windows. To fix the issue you need to switch to use of the 64-bit versions of Terraform and ZPA Terraform provider.

### Error: failed configuring the provided

This kind of error happens when the administrator fails to configure the ZPA API credentials via one of the accepted methods such as environment variables, hard-coded method (which is discouraged) or via the `credentials.json` file.

│   with provider["registry.terraform.io/zscaler/zpa"],
│   on zpa_app_connector_group.tf line 10, in provider "zpa":
│   10: provider "zpa" {}
│
│ error:Could not open credentials file, needs to contain one json object with keys: zpa_client_id, zpa_client_secret, zpa_customer_id, and
│ zpa_cloud. open /Users/wguilherme/.zpa/credentials.json: no such file or directory

### Configuration drifts with `zpa_application_segment`

The attribute `domain_names` values must be set always in lowercase. When values are set in upper case, the ZPA API automatically converts the response to lowercase which causes a drift.

To prevent that you have 2 options:

1. Set all `domain_name` values in lower case
2. Use the HCL function [lower](https://developer.hashicorp.com/terraform/language/functions/lower) to convert all cased letters in the given string to lowercase prior.

For example

```hcl
resource "zpa_application_segment" "this" {
  name              = var.name
  description       = var.description
  enabled           = var.enabled
  health_reporting  = var.health_reporting
  bypass_type       = var.bypass_type
  is_cname_enabled  = var.is_cname_enabled
  tcp_port_range    = var.tcp_port_ranges
  udp_port_range    = var.udp_port_ranges
  domain_names      = [for names in var.domain_names : lower (names)]
  segment_group_id  = var.segment_group_id
  tcp_keep_alive    = var. tcp_keep_alive
  icmp_access_type  = var.icmp_access_type
  server_groups {
    id = var.server_groups
  }
}
```

### Access Policy Operand Validation Error `zpa_policy_access_rule` Error: [WARN] when operand object type is COUNTRY_CODE LHS

This type of error happens when the administrator fails to provide a valid attribute value within the conditions.operands block of an access policy.
The below example, is an error displayed when using an invalid 2 letter `COUNTRY_CODE` wihin the policy.

│ Error: [WARN] when operand object type is COUNTRY_CODE LHS must be "valid ISO-3166 Alpha-2 country code. Please visit the following site for reference: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes" value is "00", <nil>
