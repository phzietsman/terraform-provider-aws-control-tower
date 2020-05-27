package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/phzietsman/terraform-provider-aws-control-tower/aws"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aws.Provider})
}
