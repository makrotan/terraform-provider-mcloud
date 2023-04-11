resource "mcloud_ip_block" "private" {
  name = "private"
}

resource "mcloud_ip" "mine" {
  name = "mine"
  block_id = mcloud_ip_block.private.id
  ip = "109.109.6.152"
}

resource "mcloud_ip_scope" "vault" {
  name = "vault"
}
