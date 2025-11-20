// Copyright (c) SecurityGeekIO, Inc.
// SPDX-License-Identifier: MPL-2.0

package mux

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework"
)

func MuxServer(ctx context.Context, v string) (tfprotov5.ProviderServer, error) {
	providers := []func() tfprotov5.ProviderServer{
		providerserver.NewProtocol5(framework.New(v)),
	}

	return tf5muxserver.NewMuxServer(ctx, providers...)
}
