---
page_title: "mcloud_cockroachdb Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  CockroachDB Instance.
---

# Resource `mcloud_cockroachdb`

CockroachDB Instance.



## Argument Reference

The following arguments are supported:

- `consul_cluster_id` - (Required) [string] 
- `name` - (Required) [string] 
- `version` - (Required) [string] 
- `pki_ca_id` - (Required) [string] 
- `server_pool_id` - (Required) [string] 
- `status` - [string] 
- `firewall_whitelist_ipv4` - (Required) [string] 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `access_key_primary` - [string] 
- `access_key_secondary` - [string] 
- `master_domain` - [string] 
- `ui_basic_auth_user` - [string] 
- `ui_basic_auth_password` - [string] 