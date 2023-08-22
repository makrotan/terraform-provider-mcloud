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

resource "mcloud_server_pool_hcloud" "gitlab" {
  name = "gitlab-test"
  instance_type = "cx31"
  location = "nbg1"
  instance_count = 1
  description = "Gitlab Server Pool"
}

resource "mcloud_gitlab" "gitlab" {
  name = "gtlab-test"
  fqdn = "gitlab-test.makrotan.com"
  server_pool_id = mcloud_server_pool_hcloud.gitlab.id
  version = "16.2.3-ce.0"
}

resource "mcloud_backup_policy" "gitlab" {
  name = "gitlab-test"
  ref = "Gitlab/${mcloud_gitlab.gitlab.id}"
  schedule = "10 18 * * *"
}

#
#resource "mcloud_server_pool_hcloud" "gitlab2" {
#  name = "gitlab-2"
#  instance_type = "cx21"
#  location = "nbg1"
#  instance_count = 1
#}
#
#resource "mcloud_gitlab_runner" "runner" {
#  name = "runner-test"
#  gitlab_id = mcloud_gitlab.gitlab.id
#  server_pool_id = mcloud_server_pool_hcloud.gitlab2.id
#  version = "16.1.1"
#}

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
