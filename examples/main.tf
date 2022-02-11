terraform {
  required_providers {
    mcloud = {
      version = "0.3.1"
      source = "makrotan.com/cloud/mcloud"
    }
  }
}

variable "mcloud_username" {}
variable "mcloud_password" {}

provider "mcloud" {
  username = var.mcloud_username
  password = var.mcloud_password
  host = "https://ip.makrotan.com/"
}

resource mcloud_ssh_key "foo" {
  name = "foo1"
}

output "pubkey" {
  value = mcloud_ssh_key.foo.public_key
}
output "private_key" {
  value = mcloud_ssh_key.foo.private_key
  sensitive = true
}