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

type McloudConsulCluster struct {
    AccessKeyPrimary string `json:"access_key_primary,omitempty"`
    AccessKeySecondary string `json:"access_key_secondary,omitempty"`
    EncryptionKey string `json:"encryption_key,omitempty"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    MasterDomain string `json:"master_domain,omitempty"`
    MasterServerPoolId string `json:"master_server_pool_id"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    Status string `json:"status"`
    UiBasicAuthPassword string `json:"ui_basic_auth_password,omitempty"`
    UiBasicAuthUser string `json:"ui_basic_auth_user,omitempty"`
    Version string `json:"version"`
}

type McloudConsulClusterResponse struct {
    AccessKeyPrimary string `json:"access_key_primary"`
    AccessKeySecondary string `json:"access_key_secondary"`
    EncryptionKey string `json:"encryption_key"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    MasterDomain string `json:"master_domain"`
    MasterServerPoolId string `json:"master_server_pool_id"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    Status string `json:"status"`
    UiBasicAuthPassword string `json:"ui_basic_auth_password"`
    UiBasicAuthUser string `json:"ui_basic_auth_user"`
    Version string `json:"version"`
}

func resourceMcloudConsulCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudConsulClusterCreate,
		ReadContext:   resourceMcloudConsulClusterRead,
		UpdateContext: resourceMcloudConsulClusterUpdate,
		DeleteContext: resourceMcloudConsulClusterDelete,
		Schema: map[string]*schema.Schema{
			"access_key_primary": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"access_key_secondary": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
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
			"master_domain": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"master_server_pool_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
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
			"status": &schema.Schema{
                Type:     schema.TypeString,
                Default: "running",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"ui_basic_auth_password": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"ui_basic_auth_user": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
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

func resourceMcloudConsulClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudConsulCluster{
        FirewallWhitelistIpv4: d.Get("firewall_whitelist_ipv4").(string),
        MasterServerPoolId: d.Get("master_server_pool_id").(string),
        Name: d.Get("name").(string),
        PkiCaId: d.Get("pki_ca_id").(string),
        Status: d.Get("status").(string),
        Version: d.Get("version").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/consul-cluster/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudConsulClusterResponse McloudConsulClusterResponse
	err = json.Unmarshal(body, &mcloudConsulClusterResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("access_key_primary", mcloudConsulClusterResponse.AccessKeyPrimary)
    d.Set("access_key_secondary", mcloudConsulClusterResponse.AccessKeySecondary)
    d.Set("encryption_key", mcloudConsulClusterResponse.EncryptionKey)
    d.Set("firewall_whitelist_ipv4", mcloudConsulClusterResponse.FirewallWhitelistIpv4)
    d.Set("master_domain", mcloudConsulClusterResponse.MasterDomain)
    d.Set("master_server_pool_id", mcloudConsulClusterResponse.MasterServerPoolId)
    d.Set("name", mcloudConsulClusterResponse.Name)
    d.Set("pki_ca_id", mcloudConsulClusterResponse.PkiCaId)
    d.Set("status", mcloudConsulClusterResponse.Status)
    d.Set("ui_basic_auth_password", mcloudConsulClusterResponse.UiBasicAuthPassword)
    d.Set("ui_basic_auth_user", mcloudConsulClusterResponse.UiBasicAuthUser)
    d.Set("version", mcloudConsulClusterResponse.Version)

	return diags
}

func resourceMcloudConsulClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/consul-cluster/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_consul_cluster %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudConsulClusterResponse McloudConsulClusterResponse
	err = json.Unmarshal(body, &mcloudConsulClusterResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudConsulClusterResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("access_key_primary", mcloudConsulClusterResponse.AccessKeyPrimary)
    d.Set("access_key_secondary", mcloudConsulClusterResponse.AccessKeySecondary)
    d.Set("encryption_key", mcloudConsulClusterResponse.EncryptionKey)
    d.Set("firewall_whitelist_ipv4", mcloudConsulClusterResponse.FirewallWhitelistIpv4)
    d.Set("master_domain", mcloudConsulClusterResponse.MasterDomain)
    d.Set("master_server_pool_id", mcloudConsulClusterResponse.MasterServerPoolId)
    d.Set("name", mcloudConsulClusterResponse.Name)
    d.Set("pki_ca_id", mcloudConsulClusterResponse.PkiCaId)
    d.Set("status", mcloudConsulClusterResponse.Status)
    d.Set("ui_basic_auth_password", mcloudConsulClusterResponse.UiBasicAuthPassword)
    d.Set("ui_basic_auth_user", mcloudConsulClusterResponse.UiBasicAuthUser)
    d.Set("version", mcloudConsulClusterResponse.Version)

	return diags
}
func resourceMcloudConsulClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudConsulClusterCreate(ctx, d, m)
}

func resourceMcloudConsulClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/consul-cluster/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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