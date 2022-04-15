# Server Group Example

This example will show you how to create a Server Group in the ZPA portal.
This example codifies [this API](https://help.zscaler.com/zpa/api-reference#/server-group-controller).

To run, configure your ZPA provider as described [Here](https://github.com/willguibr/terraform-provider-zpa/blob/master/docs/index.html.markdown)

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
