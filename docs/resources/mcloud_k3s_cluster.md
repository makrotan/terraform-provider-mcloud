---
page_title: "mcloud_k3s_cluster Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  K3s Kubernetes Cluster.
---

# Resource `mcloud_k3s_cluster`

K3s Kubernetes Cluster.



## Argument Reference

The following arguments are supported:

- `name` - (Required) [string]  
- `sku` - (Required) [string]  
- `version` - (Required) [string]  
- `firewall_whitelist_ipv4` - [string]  
- `master_server_pool_id` - (Required) [string]  
- `status` - [string]   (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `access_key_primary` - [string] 
- `k3s_config_yaml` - [string] 