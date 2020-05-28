# Terraform Provider for AWS ControlTower

AWS ControlTower uses Service Catalog products to provision new accounts. The official Terraform AWS Provider does not have good support for Service Catalog and this provider tries to fill the gaps. It only implements the functions needed to vend accounts using IaC.

The `controltower_account_vending` can be used to provision any Service Catalog Product, but it does try and find the vended account id from the Record Outputs as this is usually needed in other steps.

**Note** There are not tests for this provider and I may or may not end up creating some.

# Using the Provider

To use a custom-built provider in your Terraform environment, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) After placing the custom-built provider into your plugins directory,  run `terraform init` to initialize it. To get the latest release go the [Github Releases Page](https://github.com/phzietsman/terraform-provider-controltower/releases). 

`terraform-provider-controltower` implements the official AWS provider and therefore using the provider should feel very familiar. You can have all the option avaible on [AWS Provider](https://www.terraform.io/docs/providers/aws/index.html) 

```hcl
provider controltower {
  region = "eu-west-1"
}
```

---
## Data
---
## controltower_account_vending_version_lookup

TODO : Lookup the latest product artefact. 


---
## Resource
---
## controltower_account_vending

The Account Vending Resource allows the create and delete ControlTower accounts from Service Catalog. I have no idea what will happen when you update the paramters of an account vended through Service Catalog, there for I've not implemented the update method in the provider.

If the account was provisioned, the account will have to be manually closed, AWS does not support the programatic closure of accounts and the email address used for the account will be locked up of a period.

If the account was not provisioned and the Service Catalog process failed, the account can be deleted without worrying about the email address used.


### Example Usage

```hcl

variable account {
  default = "controltower-provider-1"
}

resource controltower_account_vending newaccount {

  product_id  = "prod-nl7pbqs2n3rjy"
  artefact_id = "pa-htxzmae7h7bd2"

  parameters = {
    SSOUserFirstName          = "RootName"
    SSOUserLastName           = "RootSurname"
    SSOUserEmail              = "email+${var.account}-SSO@gmail.com"
    AccountEmail              = "email+${var.account}-Account@gmail.com"
    ManagedOrganizationalUnit = "Custom"
    AccountName               = var.account

  }

  name = var.account

  # This will prevent things going sideways if you dynamically 
  # lookup the latest artefact id
  lifecycle {
    ignore_changes = [
      artefact_id,
    ]
  }
}

output account_id {
  value = controltower_account_vending.crypto.account_id
}
```

### Argument Reference

The following arguments are supported:

* `name`        - (Required) The name that will be given to the Provisioned Product
* `product_id`  - (Required) The Service Catalog Product Id
* `artefact_id` - (Required) The Product Artefact to use to create the account. 
* `parameters`  - (Required) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the Provisioned Product
* `arn` - The arn of the Provisioned Product
* `record_id` - The record of the provisioning event. I'm still figuring out how I want to handle changes to the record id

## Import

Accounts can be imported using the Provisioned Product Id, e.g.

```
$ terraform import controltower_account_vending.newaccount {{provisioned_product_id}}
```

# Contributing

This provider is based on the [AWS Provider](https://github.com/terraform-providers/terraform-provider-aws) and most those docs will apply here and will probably be more complete.


Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.10+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

Developing the Provider
---------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please check the [requirements](https://github.com/phzietsman/terraform-provider-controltower#requirements) before proceeding).

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e `$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
$ git clone git@github.com:phzietsman/terraform-provider-controltower
...
```

Enter the provider directory and run `make tools`. This will install the needed tools for the provider.

```sh
$ make tools
```

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-controltower
...
```


Testing the Provider
---------------------------

In order to test the provider, you can run `make test`.

*Note:* Make sure no `AWS_ACCESS_KEY_ID` or `AWS_SECRET_ACCESS_KEY` variables are set, and there's no `[default]` section in the AWS credentials file `~/.aws/credentials`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run. Please read [Running an Acceptance Test](https://github.com/terraform-providers/terraform-provider-aws/blob/master/.github/CONTRIBUTING.md#running-an-acceptance-test) in the contribution guidelines for more information on usage.

```sh
$ make testacc
```
