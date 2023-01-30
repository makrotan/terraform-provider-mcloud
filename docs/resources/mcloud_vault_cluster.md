---
page_title: "mcloud_vault_cluster Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  Vault Cluster.
---

# Resource `mcloud_vault_cluster`

Vault Cluster.



## Argument Reference

The following arguments are supported:

- `name` - (Required) [string]  
- `master_server_pool_id` - (Required) [string]  
- `status` - [string]   (default: `running`)
- `version` - (Required) [string]  
- `firewall_whitelist_ipv4` - (Required) [string]  
- `pki_ca_id` - (Required) [string]  

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `access_key_primary` - [string] 
- `access_key_secondary` - [string] 
- `master_domain` - [string] 
- `ui_basic_auth_user` - [string] 
- `ui_basic_auth_password` - [string] 