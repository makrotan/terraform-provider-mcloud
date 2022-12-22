---
page_title: "mcloud_mattermost Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  Mattermost Instance.
---

# Resource `mcloud_mattermost`

Mattermost Instance.



## Argument Reference

The following arguments are supported:

- `name` - (Required) [string] 
- `fqdn` - (Required) [string] 
- `sku` - (Required) [string] 
- `version` - (Required) [string] 
- `server_pool_id` - (Required) [string] 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `status` - [string] 