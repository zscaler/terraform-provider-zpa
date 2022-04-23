package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/willguibr/terraform-provider-zpa/zpa"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: zpa.Provider,
	})
}
