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

- `name` - (Required) [string]  
- `master_server_pool_id` - (Required) [string]  
- `status` - [string]   (default: `running`)
- `version` - (Required) [string]  
- `firewall_whitelist_ipv4` - (Required) [string]  
- `pki_ca_id` - (Required) [string]  
- `consul_cluster_id` - (Required) [string]  
- `vault_cluster_id` - (Required) [string]  

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `encryption_key` - [string] 
- `access_key_primary` - [string] 
- `access_key_secondary` - [string] 
- `master_domain` - [string] 
- `ui_basic_auth_user` - [string] 
- `ui_basic_auth_password` - [string] 
- `admin_ca` - [string] 
- `admin_cert` - [string] 
- `admin_key` - [string] 