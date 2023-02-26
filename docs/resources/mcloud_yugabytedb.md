---
page_title: "mcloud_yugabytedb Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  YugabyteDB Instance.
---

# Resource `mcloud_yugabytedb`

YugabyteDB Instance.



## Argument Reference

The following arguments are supported:

- `name` - (Required) [string]  
- `version` - (Required) [string]  
- `pki_ca_id` - (Required) [string]  
- `server_pool_id` - (Required) [string]  
- `status` - [string]   (default: `running`)
- `firewall_whitelist_ipv4` - [string]  

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `master_domain` - [string] 