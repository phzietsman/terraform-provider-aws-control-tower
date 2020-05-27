package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/phzietsman/terraform-provider-controltower/aws"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aws.Provider})
}
