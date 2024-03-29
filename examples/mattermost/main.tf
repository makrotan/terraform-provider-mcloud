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
  name = "mattermost-test"
  instance_type = "cpx11"
  instance_count = 1
  location = "spread"
  consul_cluster_id = mcloud_consul_cluster.foo.id
}

resource "mcloud_pki_ca" "test" {
  name = "mattermost-test"
}

resource "mcloud_mattermost" "test" {
  name           = "test"
  server_pool_id = mcloud_server_pool_hcloud.test.id
  fqdn           = "mattermost-test.makrotan.com"
  sku            = "dev"
  version        = "7.7"

  postgres_username = "foo"
  postgres_password = "bar"
}

output "out" {
  value = <<EOT
Resources successfully installed:

    Mattermost
        Access: https://${mcloud_mattermost.test.fqdn}

EOT
  sensitive = true
}
