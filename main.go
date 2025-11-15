package main

import (
	"context"
	"flag"
	"log"

	framework "github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/version"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "Start provider in debug mode.")
	flag.Parse()

	logFlags := log.Flags()
	logFlags = logFlags &^ (log.Ldate | log.Ltime)
	log.SetFlags(logFlags)

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/zscaler/zpa",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), func() provider.Provider {
		return framework.New(version.ProviderVersion)
	}, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
