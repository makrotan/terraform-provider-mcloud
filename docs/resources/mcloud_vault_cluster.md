---
page_title: "mcloud_vault_cluster Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  VaultCluster(status, backup_ref, name, created, master_server_pool, version, access_key_primary, access_key_secondary, pki_ca, seal_keys, root_token, ip_scope)
---

# Resource `mcloud_vault_cluster`

VaultCluster(status, backup_ref, name, created, master_server_pool, version, access_key_primary, access_key_secondary, pki_ca, seal_keys, root_token, ip_scope)



## Argument Reference

The following arguments are supported:

- `ip_scope_id` - [string]  
- `master_server_pool_id` - (Required) [string]  
- `name` - (Required) [string]  
- `pki_ca_id` - [string]  
- `status` - [string]   (default: `running`)
- `version` - (Required) [string] Possible values: `1.12.2-1`, `1.13.1-1` 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `access_key_primary` - [string] 
- `access_key_secondary` - [string] 
- `master_domain` - [string] 
- `root_token` - [string] 
- `ui_basic_auth_password` - [string] 
- `ui_basic_auth_user` - [string] 