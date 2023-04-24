resource "mcloud_consul_cluster" "test" {
  name = "consul-test"
  master_server_pool_id = mcloud_server_pool_hcloud.vaco.id
  pki_ca_id = mcloud_pki_ca.test.id
  ip_scope_id = mcloud_ip_scope.consul.id
  version = "1.15.2"

  # we are using the same server pools for everything, so we need to take care of not doing things concurrently
  depends_on = [
    mcloud_vault_cluster.test,
  ]
}
