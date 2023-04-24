
resource "mcloud_server_pool_hcloud" "client_a" {
  name = "nomad-client-A"
  instance_type = "cpx11"
  instance_count = 1
  location = "spread"
  consul_cluster_id = mcloud_consul_cluster.test.id
}

resource "mcloud_nomad_server_pool" "client_a" {
  name = "nomad-client-A"
  nomad_cluster_id = mcloud_nomad_cluster.test.id
  server_pool_id = mcloud_server_pool_hcloud.client_a.id
}
