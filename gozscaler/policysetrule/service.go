package policysetrule

import (
	"github.com/willguibr/terraform-provider-zpa/gozscaler"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
)

type Service struct {
	Client *client.Client
}

func New(c *gozscaler.Config) *Service {
	return &Service{Client: client.NewClient(c)}
}
