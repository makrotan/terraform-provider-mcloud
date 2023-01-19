package mcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type McloudK3sCluster struct {
    Name string `json:"name"`
    Sku string `json:"sku"`
    Version string `json:"version"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    MasterServerPoolId string `json:"master_server_pool_id"`
    Status string `json:"status"`
    AccessKeyPrimary string `json:"access_key_primary,omitempty"`
    K3sConfigYaml string `json:"k3s_config_yaml,omitempty"`
}

type McloudK3sClusterResponse struct {
    Name string `json:"name"`
    Sku string `json:"sku"`
    Version string `json:"version"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    MasterServerPoolId string `json:"master_server_pool_id"`
    Status string `json:"status"`
    AccessKeyPrimary string `json:"access_key_primary"`
    K3sConfigYaml string `json:"k3s_config_yaml"`
}

func resourceMcloudK3sCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudK3sClusterCreate,
		ReadContext:   resourceMcloudK3sClusterRead,
		UpdateContext: resourceMcloudK3sClusterUpdate,
		DeleteContext: resourceMcloudK3sClusterDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"sku": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: true,
			},
			"version": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"firewall_whitelist_ipv4": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"master_server_pool_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: true,
			},
			"status": &schema.Schema{
                Type:     schema.TypeString,
                Default: "running",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"access_key_primary": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"k3s_config_yaml": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMcloudK3sClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudK3sCluster{
        Name: d.Get("name").(string),
        Sku: d.Get("sku").(string),
        Version: d.Get("version").(string),
        FirewallWhitelistIpv4: d.Get("firewall_whitelist_ipv4").(string),
        MasterServerPoolId: d.Get("master_server_pool_id").(string),
        Status: d.Get("status").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/k3s-cluster/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.api_token))

	res, err := provider.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudK3sClusterResponse McloudK3sClusterResponse
	err = json.Unmarshal(body, &mcloudK3sClusterResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("name", mcloudK3sClusterResponse.Name)
    d.Set("sku", mcloudK3sClusterResponse.Sku)
    d.Set("version", mcloudK3sClusterResponse.Version)
    d.Set("firewall_whitelist_ipv4", mcloudK3sClusterResponse.FirewallWhitelistIpv4)
    d.Set("master_server_pool_id", mcloudK3sClusterResponse.MasterServerPoolId)
    d.Set("status", mcloudK3sClusterResponse.Status)
    d.Set("access_key_primary", mcloudK3sClusterResponse.AccessKeyPrimary)
    d.Set("k3s_config_yaml", mcloudK3sClusterResponse.K3sConfigYaml)

	return diags
}

func resourceMcloudK3sClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/k3s-cluster/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.api_token))

	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode == 404 {
		log.Printf("[WARN] mcloud_k3s_cluster %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudK3sClusterResponse McloudK3sClusterResponse
	err = json.Unmarshal(body, &mcloudK3sClusterResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudK3sClusterResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("name", mcloudK3sClusterResponse.Name)
    d.Set("sku", mcloudK3sClusterResponse.Sku)
    d.Set("version", mcloudK3sClusterResponse.Version)
    d.Set("firewall_whitelist_ipv4", mcloudK3sClusterResponse.FirewallWhitelistIpv4)
    d.Set("master_server_pool_id", mcloudK3sClusterResponse.MasterServerPoolId)
    d.Set("status", mcloudK3sClusterResponse.Status)
    d.Set("access_key_primary", mcloudK3sClusterResponse.AccessKeyPrimary)
    d.Set("k3s_config_yaml", mcloudK3sClusterResponse.K3sConfigYaml)

	return diags
}
func resourceMcloudK3sClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudK3sClusterCreate(ctx, d, m)
}

func resourceMcloudK3sClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/k3s-cluster/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.api_token))

	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode >= 300 {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	return diags
}