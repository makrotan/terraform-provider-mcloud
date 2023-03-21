---
page_title: "mcloud_consul_cluster Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  Consul Cluster.
---

# Resource `mcloud_consul_cluster`

Consul Cluster.



## Argument Reference

The following arguments are supported:

- `firewall_whitelist_ipv4` - (Required) [string]  
- `master_server_pool_id` - (Required) [string]  
- `name` - (Required) [string]  
- `pki_ca_id` - (Required) [string]  
- `status` - [string]   (default: `running`)
- `version` - (Required) [string]  

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `access_key_primary` - [string] 
- `access_key_secondary` - [string] 
- `encryption_key` - [string] 
- `master_domain` - [string] 
- `ui_basic_auth_password` - [string] 
- `ui_basic_auth_user` - [string] 