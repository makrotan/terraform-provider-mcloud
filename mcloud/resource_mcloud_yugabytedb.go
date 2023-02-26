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

type McloudYugabytedb struct {
    Name string `json:"name"`
    Version string `json:"version"`
    PkiCaId string `json:"pki_ca_id"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    MasterDomain string `json:"master_domain,omitempty"`
}

type McloudYugabytedbResponse struct {
    Name string `json:"name"`
    Version string `json:"version"`
    PkiCaId string `json:"pki_ca_id"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    MasterDomain string `json:"master_domain"`
}

func resourceMcloudYugabytedb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudYugabytedbCreate,
		ReadContext:   resourceMcloudYugabytedbRead,
		UpdateContext: resourceMcloudYugabytedbUpdate,
		DeleteContext: resourceMcloudYugabytedbDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"version": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"pki_ca_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: true,
			},
			"server_pool_id": &schema.Schema{
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
			"firewall_whitelist_ipv4": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"master_domain": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMcloudYugabytedbCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudYugabytedb{
        Name: d.Get("name").(string),
        Version: d.Get("version").(string),
        PkiCaId: d.Get("pki_ca_id").(string),
        ServerPoolId: d.Get("server_pool_id").(string),
        Status: d.Get("status").(string),
        FirewallWhitelistIpv4: d.Get("firewall_whitelist_ipv4").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/yugabytedb/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudYugabytedbResponse McloudYugabytedbResponse
	err = json.Unmarshal(body, &mcloudYugabytedbResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("name", mcloudYugabytedbResponse.Name)
    d.Set("version", mcloudYugabytedbResponse.Version)
    d.Set("pki_ca_id", mcloudYugabytedbResponse.PkiCaId)
    d.Set("server_pool_id", mcloudYugabytedbResponse.ServerPoolId)
    d.Set("status", mcloudYugabytedbResponse.Status)
    d.Set("firewall_whitelist_ipv4", mcloudYugabytedbResponse.FirewallWhitelistIpv4)
    d.Set("master_domain", mcloudYugabytedbResponse.MasterDomain)

	return diags
}

func resourceMcloudYugabytedbRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/yugabytedb/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_yugabytedb %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudYugabytedbResponse McloudYugabytedbResponse
	err = json.Unmarshal(body, &mcloudYugabytedbResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudYugabytedbResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("name", mcloudYugabytedbResponse.Name)
    d.Set("version", mcloudYugabytedbResponse.Version)
    d.Set("pki_ca_id", mcloudYugabytedbResponse.PkiCaId)
    d.Set("server_pool_id", mcloudYugabytedbResponse.ServerPoolId)
    d.Set("status", mcloudYugabytedbResponse.Status)
    d.Set("firewall_whitelist_ipv4", mcloudYugabytedbResponse.FirewallWhitelistIpv4)
    d.Set("master_domain", mcloudYugabytedbResponse.MasterDomain)

	return diags
}
func resourceMcloudYugabytedbUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudYugabytedbCreate(ctx, d, m)
}

func resourceMcloudYugabytedbDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/yugabytedb/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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