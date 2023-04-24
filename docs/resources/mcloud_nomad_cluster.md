---
page_title: "mcloud_nomad_cluster Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  NomadCluster(status, backup_ref, name, created, master_server_pool, encryption_key, version, ip_scope, access_key_primary, access_key_secondary, pki_ca, consul_cluster, vault_cluster, vault_token)
---

# Resource `mcloud_nomad_cluster`

NomadCluster(status, backup_ref, name, created, master_server_pool, encryption_key, version, ip_scope, access_key_primary, access_key_secondary, pki_ca, consul_cluster, vault_cluster, vault_token)



## Argument Reference

The following arguments are supported:

- `consul_cluster_id` - (Required) [string]  
- `ip_scope_id` - [string]  
- `master_server_pool_id` - [string]  
- `name` - (Required) [string]  
- `pki_ca_id` - [string]  
- `status` - [string]   (default: `running`)
- `vault_cluster_id` - [string]  
- `version` - (Required) [string] Possible values: `1.2.7-1`, `1.3.0-1`, `1.3.8-1`, `1.4.3-1`, `1.5.3-1` 

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
- `vault_token` - [string] 