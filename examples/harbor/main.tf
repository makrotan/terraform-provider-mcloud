terraform {
  required_providers {
    mcloud = {
      version = "0.1.0"
      source = "makrotan.com/cloud/mcloud"
    }
  }
}

variable "mcloud_token" {}

provider "mcloud" {
  api_token = var.mcloud_token
  host = "http://127.0.0.1:10004/"
}

resource "mcloud_server_pool_hcloud" "foo" {
  name = "foo"
  instance_type = "cpx11"
  instance_count = 1
}
