terraform {
  required_providers {
    jsonfile = {
      source = "registry.terraform.io/papegaaij/jsonfile"
    }
  }
}

provider "jsonfile" {
}

resource "jsonfile_data" "data" {
  value = "test 2"
#  nested = [{
#    fixed = "fixed"
#  }]
}
