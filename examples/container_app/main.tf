terraform {
  required_providers {
    mcloud = {
      version = "0.2.0"
      source = "makrotan/mcloud"
    }
  }
}

variable "mcloud_token" {}

provider "mcloud" {
  api_token = var.mcloud_token
  host = "http://127.0.0.1:10004/"
}

resource "mcloud_server_pool_hcloud" "ctest1" {
  name = "container-test-1"
  instance_type = "cx21"
  location = "nbg1"
  instance_count = 1
}
resource "mcloud_server_pool_hcloud" "ctest2" {
  name = "container-test-2"
  instance_type = "cx21"
  location = "nbg1"
  instance_count = 1
}

resource "mcloud_container_app" "ghost" {
  name = "container-test"
  fqdn = "foo.leanit.gmbh"
  server_pool_id = mcloud_server_pool_hcloud.ctest1.id
  definition = file("${path.module}/docker-compose.yml")
  depends_on = [mcloud_server_pool_hcloud.ctest1]
}

resource "mcloud_container_app" "ghost2" {
  name = "ghost2"
  fqdn = "ghost2.leanit.gmbh"
  server_pool_id = mcloud_server_pool_hcloud.ctest1.id
  definition = file("${path.module}/docker-compose.yml")
  depends_on = [mcloud_server_pool_hcloud.ctest1]
}

#
#
#output "out" {
#  value = <<EOT
#Resources successfully installed:
#
#    Gitlab:
#        Access: https://${mcloud_gitlab.gitlab.fqdn}/
#        User: ${mcloud_gitlab.gitlab.admin_username}
#        Password: ${mcloud_gitlab.gitlab.admin_password}
#
#EOT
#  sensitive = true
#}
#
