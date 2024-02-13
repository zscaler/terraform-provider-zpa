# Create Policy Redirection Rule

This example will show you how to create a policy redirection rule to allow you to set criteria for preferring ZPA Private Service Edges over ZPA Public Service Edges
This example codifies [this API](https://help.zscaler.com/zpa/about-redirection-policy).

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
