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
  name = "foo"
}

#
#module "psl" {
#  source = "./coffee"
#
#  coffee_name = "Packer Spiced Latte"
#}
#
#output "psl" {
#  value = module.psl.coffee
#}
#
#data "hashicups_order" "order" {
#  id = 1
#}
#
#output "order" {
#  value = data.hashicups_order.order
#}
#
#resource "hashicups_order" "edu" {
#  items {
#    coffee {
#      id = 3
#    }
#    quantity = 3
#  }
#  items {
#    coffee {
#      id = 2
#    }
#    quantity = 3
#  }
#}
#
#output "edu_order" {
#  value = hashicups_order.edu
#}
#
#
#data "hashicups_order" "first" {
#  id = 1
#}
#
#output "first_order" {
#  value = data.hashicups_order.first
#}
