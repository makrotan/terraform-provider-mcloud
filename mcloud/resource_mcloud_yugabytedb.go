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
    ConsulClusterId string `json:"consul_cluster_id"`
    EncryptionKey string `json:"encryption_key,omitempty"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    Fqdn string `json:"fqdn,omitempty"`
    IpScopeAdminId string `json:"ip_scope_admin_id"`
    IpScopeClientId string `json:"ip_scope_client_id"`
    Meta map[string]interface{} `json:"meta"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
    Version string `json:"version"`
}

type McloudYugabytedbResponse struct {
    ConsulClusterId string `json:"consul_cluster_id"`
    EncryptionKey string `json:"encryption_key"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    Fqdn string `json:"fqdn"`
    IpScopeAdminId string `json:"ip_scope_admin_id"`
    IpScopeClientId string `json:"ip_scope_client_id"`
    Meta map[string]interface{} `json:"meta"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
    Version string `json:"version"`
}

func resourceMcloudYugabytedb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudYugabytedbCreate,
		ReadContext:   resourceMcloudYugabytedbRead,
		UpdateContext: resourceMcloudYugabytedbUpdate,
		DeleteContext: resourceMcloudYugabytedbDelete,
		Schema: map[string]*schema.Schema{
			"consul_cluster_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"encryption_key": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"firewall_whitelist_ipv4": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"fqdn": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"ip_scope_admin_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"ip_scope_client_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"meta": &schema.Schema{
                Type:     schema.TypeMap,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"pki_ca_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"server_pool_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"status": &schema.Schema{
                Type:     schema.TypeString,
                Default: "running",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"version": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
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
        ConsulClusterId: d.Get("consul_cluster_id").(string),
        FirewallWhitelistIpv4: d.Get("firewall_whitelist_ipv4").(string),
        IpScopeAdminId: d.Get("ip_scope_admin_id").(string),
        IpScopeClientId: d.Get("ip_scope_client_id").(string),
        Meta: d.Get("meta").(map[string]interface{}),
        Name: d.Get("name").(string),
        PkiCaId: d.Get("pki_ca_id").(string),
        ServerPoolId: d.Get("server_pool_id").(string),
        Status: d.Get("status").(string),
        Version: d.Get("version").(string),
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
    d.Set("consul_cluster_id", mcloudYugabytedbResponse.ConsulClusterId)
    d.Set("encryption_key", mcloudYugabytedbResponse.EncryptionKey)
    d.Set("firewall_whitelist_ipv4", mcloudYugabytedbResponse.FirewallWhitelistIpv4)
    d.Set("fqdn", mcloudYugabytedbResponse.Fqdn)
    d.Set("ip_scope_admin_id", mcloudYugabytedbResponse.IpScopeAdminId)
    d.Set("ip_scope_client_id", mcloudYugabytedbResponse.IpScopeClientId)
    d.Set("meta", mcloudYugabytedbResponse.Meta)
    d.Set("name", mcloudYugabytedbResponse.Name)
    d.Set("pki_ca_id", mcloudYugabytedbResponse.PkiCaId)
    d.Set("server_pool_id", mcloudYugabytedbResponse.ServerPoolId)
    d.Set("status", mcloudYugabytedbResponse.Status)
    d.Set("version", mcloudYugabytedbResponse.Version)

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
    d.Set("consul_cluster_id", mcloudYugabytedbResponse.ConsulClusterId)
    d.Set("encryption_key", mcloudYugabytedbResponse.EncryptionKey)
    d.Set("firewall_whitelist_ipv4", mcloudYugabytedbResponse.FirewallWhitelistIpv4)
    d.Set("fqdn", mcloudYugabytedbResponse.Fqdn)
    d.Set("ip_scope_admin_id", mcloudYugabytedbResponse.IpScopeAdminId)
    d.Set("ip_scope_client_id", mcloudYugabytedbResponse.IpScopeClientId)
    d.Set("meta", mcloudYugabytedbResponse.Meta)
    d.Set("name", mcloudYugabytedbResponse.Name)
    d.Set("pki_ca_id", mcloudYugabytedbResponse.PkiCaId)
    d.Set("server_pool_id", mcloudYugabytedbResponse.ServerPoolId)
    d.Set("status", mcloudYugabytedbResponse.Status)
    d.Set("version", mcloudYugabytedbResponse.Version)

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