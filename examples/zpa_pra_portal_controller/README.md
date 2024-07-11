
# Create Privileged Remote Access Portal

This example will show you how to create a privileged remote access portal.
This example codifies [this API](https://help.zscaler.com/zpa/configuring-privileged-portals-using-api).

To run, configure your ZPA provider as described [Here](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/index.md)

## Run the example

From inside of this directory:

```bash
terraform init
terraform plan -out theplan
terraform apply theplan
```

## Destroy 💥

```bash
terraform destroy
```
