---
page_title: "mcloud_nomad_cluster Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  Nomad Cluster.
---

# Resource `mcloud_nomad_cluster`

Nomad Cluster.



## Argument Reference

The following arguments are supported:

- `consul_cluster_id` - (Required) [string]  
- `firewall_whitelist_ipv4` - (Required) [string]  
- `master_server_pool_id` - (Required) [string]  
- `name` - (Required) [string]  
- `pki_ca_id` - (Required) [string]  
- `status` - [string]   (default: `running`)
- `vault_cluster_id` - (Required) [string]  
- `version` - (Required) [string]  

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `access_key_primary` - [string] 
- `access_key_secondary` - [string] 
- `admin_ca` - [string] 
- `admin_cert` - [string] 
- `admin_key` - [string] 
- `encryption_key` - [string] 
- `master_domain` - [string] 
- `ui_basic_auth_password` - [string] 
- `ui_basic_auth_user` - [string] 