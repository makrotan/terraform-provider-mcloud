---
page_title: "mcloud_pki_ca Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  PKI Certificate Authority.
---

# Resource `mcloud_pki_ca`

PKI Certificate Authority.



## Argument Reference

The following arguments are supported:

- `name` - (Required) [string]  
- `valid_days` - [number]   (default: `7300`)
- `country` - [string]   (default: `DE`)
- `state` - [string]   (default: `DE`)
- `city` - [string]   (default: `Leipzig`)
- `organisation` - [string]   (default: `Makrotan`)
- `unit` - [string]   (default: `IT`)
- `email` - [string]   (default: `info@makrotan.com`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `key_pub` - [string] 
- `key_priv` - [string] 