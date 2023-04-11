---
page_title: "mcloud_yugabytedb Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  Yugabytedb(status, backup_ref, name, created, server_pool, firewall_whitelist_ipv4, meta, encryption_key, version, pki_ca, ip_scope_admin, ip_scope_client, consul_cluster)
---

# Resource `mcloud_yugabytedb`

Yugabytedb(status, backup_ref, name, created, server_pool, firewall_whitelist_ipv4, meta, encryption_key, version, pki_ca, ip_scope_admin, ip_scope_client, consul_cluster)



## Argument Reference

The following arguments are supported:

- `consul_cluster_id` - [string]  
- `firewall_whitelist_ipv4` - [string]  
- `ip_scope_admin_id` - [string]  
- `ip_scope_client_id` - [string]  
- `meta` - [map]  
- `name` - (Required) [string]  
- `pki_ca_id` - [string]  
- `server_pool_id` - [string]  
- `status` - [string]   (default: `running`)
- `version` - (Required) [string] Possible values: `2.17.1.0` 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `encryption_key` - [string] 
- `fqdn` - [string] 