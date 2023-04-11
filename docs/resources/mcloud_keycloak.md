---
page_title: "mcloud_keycloak Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  Keycloak(status, backup_ref, name, created, server_pool, admin_password, fqdn, pki_ca, app_port, secret_key, sku, version, themes)
---

# Resource `mcloud_keycloak`

Keycloak(status, backup_ref, name, created, server_pool, admin_password, fqdn, pki_ca, app_port, secret_key, sku, version, themes)



## Argument Reference

The following arguments are supported:

- `fqdn` - (Required) [string]  
- `name` - (Required) [string]  
- `pki_ca_id` - [string]  
- `server_pool_id` - [string]  
- `sku` - (Required) [string] Possible values: `dev` 
- `status` - [string]   (default: `running`)
- `themes` - [string]  
- `version` - (Required) [string] Possible values: `21.0.1` 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `admin_password` - [string] 
- `secret_key` - [string] 