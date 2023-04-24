#resource "mcloud_server_pool_hcloud" "foo" {
#  name = "foo-test"
#  instance_type = "cpx11"
#  instance_count = 1
#  location = "spread"
#}
#
#resource "mcloud_ip_scope" "foo" {
#  name = "foo"
#}
#
#resource "mcloud_pki_ca" "foo" {
#  name = "foo"
#}
#
#resource "mcloud_consul_cluster" "foo" {
#  name = "consul-foo"
#  master_server_pool_id = mcloud_server_pool_hcloud.foo.id
#  pki_ca_id = mcloud_pki_ca.foo.id
#  ip_scope_id = mcloud_ip_scope.foo.id
#  version = "1.15.2"
#}
