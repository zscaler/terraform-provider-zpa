package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/zscaler/terraform-provider-zpa/zpa"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: zpa.Provider,
	})
}
