---
page_title: "mcloud_consul_cluster Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  ConsulCluster(name, status, backup_ref, created, master_server_pool, encryption_key, version, ip_scope, access_key_primary, access_key_secondary, pki_ca)
---

# Resource `mcloud_consul_cluster`

ConsulCluster(name, status, backup_ref, created, master_server_pool, encryption_key, version, ip_scope, access_key_primary, access_key_secondary, pki_ca)



## Argument Reference

The following arguments are supported:

- `ip_scope_id` - [string]  
- `master_server_pool_id` - (Required) [string]  
- `name` - (Required) [string]  
- `pki_ca_id` - [string]  
- `status` - [string]   (default: `running`)
- `version` - (Required) [string] Possible values: `1.11.4`, `1.11.5`, `1.12.0`, `1.14.3`, `1.15.2` 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `access_key_primary` - [string] 
- `access_key_secondary` - [string] 
- `encryption_key` - [string] 
- `master_domain` - [string] 
- `ui_basic_auth_password` - [string] 
- `ui_basic_auth_user` - [string] 