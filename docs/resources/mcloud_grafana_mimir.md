---
page_title: "mcloud_grafana_mimir Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  GrafanaMimir(backup_ref, name, created, server_pool, access_key, fqdn, app_port, pki_ca, status, version)
---

# Resource `mcloud_grafana_mimir`

GrafanaMimir(backup_ref, name, created, server_pool, access_key, fqdn, app_port, pki_ca, status, version)



## Argument Reference

The following arguments are supported:

- `fqdn` - (Required) [string]  
- `name` - (Required) [string]  
- `pki_ca_id` - [string]  
- `server_pool_id` - [string]  
- `status` - [string]   (default: `running`)
- `version` - (Required) [string] Possible values: `2.6.0` 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `access_key` - [string] 
- `basic_auth_password` - [string] 
- `basic_auth_user` - [string] 