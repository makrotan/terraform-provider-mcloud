---
page_title: "mcloud_ssh_key Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  SSH Key.
---

# Resource `mcloud_ssh_key`

SSH Key.



## Argument Reference

The following arguments are supported:

- `name` - (Required) [string]  
- `private_key` - (Required) [string] private key in pem format: `-----BEGIN OPENSSH PRIVATE KEY-----\nb3Bl...cgECAw==\n-----END OPENSSH PRIVATE KEY-----\n` 
- `public_key` - (Required) [string] public key as givin in an authorized_hosts file, `ssh-rsa AAAAB3Nz...PCmXzzFLKoC0Agvc= hostname` 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:
