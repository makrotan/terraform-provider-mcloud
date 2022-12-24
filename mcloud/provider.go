package mcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"api_token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mcloud_consul_cluster":        resourceMcloudConsulCluster(),
			"mcloud_cockroachdb":           resourceMcloudCockroachdb(),
			"mcloud_erpnext":               resourceMcloudErpnext(),
			"mcloud_mattermost":            resourceMcloudMattermost(),
			"mcloud_pki_ca":                resourceMcloudPkiCa(),
			"mcloud_ssh_key":               resourceMcloudSshKey(),
			"mcloud_server_pool_dedicated": resourceMcloudServerPoolDedicated(),
			"mcloud_server_pool_hcloud":    resourceMcloudServerPoolHcloud(),
			"mcloud_server_dedicated":      resourceMcloudServerDedicated(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var host *string

	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	api_token := d.Get("api_token").(string)
	c, err := NewClient(host, nil, nil, &api_token)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create mcloud client",
			Detail:   "Unable to authenticate mcloud client using bearer_token: " + err.Error(),
		})

		return nil, diags
	}
	c.api_token = d.Get("api_token").(string)

	return c, diags
}
