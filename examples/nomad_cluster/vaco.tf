resource "mcloud_server_pool_hcloud" "vaco" {
  name = "nomad-vaco"
  instance_type = "cpx11"
  instance_count = 1
  location = "spread"
}

