---
page_title: "mcloud_pki_cert Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  PKI Certificate.
---

# Resource `mcloud_pki_cert`

PKI Certificate.



## Argument Reference

The following arguments are supported:

- `ca_id` - (Required) [string]  
- `common_name` - (Required) [string]  
- `name` - (Required) [string]  

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `key_priv` - [string] 
- `key_pub` - [string] 