provider controltower {
  region = "eu-west-1"
}

variable account {
  default = "controltower-provider-1"
}

resource controltower_account_vending crypto {

  product_id  = "prod-nl7pbqs2n3rjy"
  artefact_id = "pa-htxzmae7h7bd2"

  parameters = {
    SSOUserFirstName          = "RootName"
    SSOUserLastName           = "RootSurname"
    SSOUserEmail              = "phzietsman+${var.account}-SSO@gmail.com"
    AccountEmail              = "phzietsman+${var.account}-Account@gmail.com"
    ManagedOrganizationalUnit = "Custom"
    AccountName               = var.account

  }

  name = var.account
}

output account_id {
  value = controltower_account_vending.crypto.account_id
}