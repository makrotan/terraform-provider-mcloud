---
page_title: "mcloud_server_dedicated Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  ServerDedicated(name, status, backup_ref, created, ipv4, ipv6, provider, provider_ref, region, az, pool, memory, cpu_cores, disk_size, price_per_month)
---

# Resource `mcloud_server_dedicated`

ServerDedicated(name, status, backup_ref, created, ipv4, ipv6, provider, provider_ref, region, az, pool, memory, cpu_cores, disk_size, price_per_month)



## Argument Reference

The following arguments are supported:

- `az` - (Required) [string] Possible values: `hetzner-fsn1`, `hetzner-nbg1`, `hetzner-hel1`, `contabo-fra`, `contabo-nbg`, `netcup-nbg`, `hoston-ffm` 
- `cpu_cores` - [number]  
- `disk_size` - [number]  
- `ipv4` - (Required) [string]  
- `ipv6` - [string]  
- `memory` - [number]  
- `name` - (Required) [string]  
- `pool_id` - [string]  
- `price_per_month` - [number]  
- `provider_id` - (Required) [string]  
- `provider_ref` - [string]  
- `region` - (Required) [string] Possible values: `europe`, `na`, `sa`, `au`, `asia` 
- `status` - [string]   (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:
