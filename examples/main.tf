terraform {
  required_providers {
    hashicups = {
      version = "0.3.1"
      source = "hashicorp.com/edu/hashicups"
    }
  }
}

variable "mcloud_username" {}
variable "mcloud_password" {}

provider "hashicups" {
  username = var.mcloud_username
  password = var.mcloud_password
  host = "https://ip.makrotan.com/"
}

resource hashicups_mcloud_ssh_key "foo" {
  name = "foo2"
}
