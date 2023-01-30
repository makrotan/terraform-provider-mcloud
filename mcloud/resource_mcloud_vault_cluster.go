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

type McloudVaultCluster struct {
    Name string `json:"name"`
    MasterServerPoolId string `json:"master_server_pool_id"`
    Status string `json:"status"`
    Version string `json:"version"`
    AccessKeyPrimary string `json:"access_key_primary,omitempty"`
    AccessKeySecondary string `json:"access_key_secondary,omitempty"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    PkiCaId string `json:"pki_ca_id"`
    MasterDomain string `json:"master_domain,omitempty"`
    UiBasicAuthUser string `json:"ui_basic_auth_user,omitempty"`
    UiBasicAuthPassword string `json:"ui_basic_auth_password,omitempty"`
}

type McloudVaultClusterResponse struct {
    Name string `json:"name"`
    MasterServerPoolId string `json:"master_server_pool_id"`
    Status string `json:"status"`
    Version string `json:"version"`
    AccessKeyPrimary string `json:"access_key_primary"`
    AccessKeySecondary string `json:"access_key_secondary"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    PkiCaId string `json:"pki_ca_id"`
    MasterDomain string `json:"master_domain"`
    UiBasicAuthUser string `json:"ui_basic_auth_user"`
    UiBasicAuthPassword string `json:"ui_basic_auth_password"`
}

func resourceMcloudVaultCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudVaultClusterCreate,
		ReadContext:   resourceMcloudVaultClusterRead,
		UpdateContext: resourceMcloudVaultClusterUpdate,
		DeleteContext: resourceMcloudVaultClusterDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
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
			"version": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"access_key_primary": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"access_key_secondary": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"firewall_whitelist_ipv4": &schema.Schema{
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
			"master_domain": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"ui_basic_auth_user": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"ui_basic_auth_password": &schema.Schema{
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

func resourceMcloudVaultClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudVaultCluster{
        Name: d.Get("name").(string),
        MasterServerPoolId: d.Get("master_server_pool_id").(string),
        Status: d.Get("status").(string),
        Version: d.Get("version").(string),
        FirewallWhitelistIpv4: d.Get("firewall_whitelist_ipv4").(string),
        PkiCaId: d.Get("pki_ca_id").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/vault-cluster/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudVaultClusterResponse McloudVaultClusterResponse
	err = json.Unmarshal(body, &mcloudVaultClusterResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("name", mcloudVaultClusterResponse.Name)
    d.Set("master_server_pool_id", mcloudVaultClusterResponse.MasterServerPoolId)
    d.Set("status", mcloudVaultClusterResponse.Status)
    d.Set("version", mcloudVaultClusterResponse.Version)
    d.Set("access_key_primary", mcloudVaultClusterResponse.AccessKeyPrimary)
    d.Set("access_key_secondary", mcloudVaultClusterResponse.AccessKeySecondary)
    d.Set("firewall_whitelist_ipv4", mcloudVaultClusterResponse.FirewallWhitelistIpv4)
    d.Set("pki_ca_id", mcloudVaultClusterResponse.PkiCaId)
    d.Set("master_domain", mcloudVaultClusterResponse.MasterDomain)
    d.Set("ui_basic_auth_user", mcloudVaultClusterResponse.UiBasicAuthUser)
    d.Set("ui_basic_auth_password", mcloudVaultClusterResponse.UiBasicAuthPassword)

	return diags
}

func resourceMcloudVaultClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/vault-cluster/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_vault_cluster %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudVaultClusterResponse McloudVaultClusterResponse
	err = json.Unmarshal(body, &mcloudVaultClusterResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudVaultClusterResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("name", mcloudVaultClusterResponse.Name)
    d.Set("master_server_pool_id", mcloudVaultClusterResponse.MasterServerPoolId)
    d.Set("status", mcloudVaultClusterResponse.Status)
    d.Set("version", mcloudVaultClusterResponse.Version)
    d.Set("access_key_primary", mcloudVaultClusterResponse.AccessKeyPrimary)
    d.Set("access_key_secondary", mcloudVaultClusterResponse.AccessKeySecondary)
    d.Set("firewall_whitelist_ipv4", mcloudVaultClusterResponse.FirewallWhitelistIpv4)
    d.Set("pki_ca_id", mcloudVaultClusterResponse.PkiCaId)
    d.Set("master_domain", mcloudVaultClusterResponse.MasterDomain)
    d.Set("ui_basic_auth_user", mcloudVaultClusterResponse.UiBasicAuthUser)
    d.Set("ui_basic_auth_password", mcloudVaultClusterResponse.UiBasicAuthPassword)

	return diags
}
func resourceMcloudVaultClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudVaultClusterCreate(ctx, d, m)
}

func resourceMcloudVaultClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/vault-cluster/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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