terraform {
  required_providers {
    mcloud = {
      version = "0.2.0"
      source = "makrotan/mcloud"
    }
  }
}

variable "mcloud_token" {}
variable "firewall_whitelist_ipv4" {
  default = ""
}

provider "mcloud" {
  api_token = var.mcloud_token
  host = "http://127.0.0.1:10004/"
}

resource "mcloud_server_pool_hcloud" "test" {
  name = "erpnext-test"
  instance_type = "cpx11"
  instance_count = 1
  location = "nbg1"
}

resource "mcloud_pki_ca" "test" {
  name = "erpnext-test"
}

resource "mcloud_erpnext" "test" {
  name           = "test"
  server_pool_id = mcloud_server_pool_hcloud.test.id
  fqdn           = "erpnext-test.makrotan.com"
  sku            = "dev"
  version = "14.13.0"
#  version        = "14.12.1"
}


output "out" {
  value = <<EOT
Resources successfully installed:

    ErpNext
        Access: https://${mcloud_erpnext.test.fqdn}
        Username: administrator
        Password: ${mcloud_erpnext.test.admin_password}

EOT
  sensitive = true
}

