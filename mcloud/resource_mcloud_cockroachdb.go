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

type McloudCockroachdb struct {
    AccessKeyPrimary string `json:"access_key_primary,omitempty"`
    AccessKeySecondary string `json:"access_key_secondary,omitempty"`
    ConsulClusterId string `json:"consul_cluster_id"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    MasterDomain string `json:"master_domain,omitempty"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
    UiBasicAuthPassword string `json:"ui_basic_auth_password,omitempty"`
    UiBasicAuthUser string `json:"ui_basic_auth_user,omitempty"`
    Version string `json:"version"`
}

type McloudCockroachdbResponse struct {
    AccessKeyPrimary string `json:"access_key_primary"`
    AccessKeySecondary string `json:"access_key_secondary"`
    ConsulClusterId string `json:"consul_cluster_id"`
    FirewallWhitelistIpv4 string `json:"firewall_whitelist_ipv4"`
    MasterDomain string `json:"master_domain"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
    UiBasicAuthPassword string `json:"ui_basic_auth_password"`
    UiBasicAuthUser string `json:"ui_basic_auth_user"`
    Version string `json:"version"`
}

func resourceMcloudCockroachdb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudCockroachdbCreate,
		ReadContext:   resourceMcloudCockroachdbRead,
		UpdateContext: resourceMcloudCockroachdbUpdate,
		DeleteContext: resourceMcloudCockroachdbDelete,
		Schema: map[string]*schema.Schema{
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
			"consul_cluster_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: true,
			},
			"firewall_whitelist_ipv4": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"master_domain": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
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
			"ui_basic_auth_password": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"ui_basic_auth_user": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
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

func resourceMcloudCockroachdbCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudCockroachdb{
        ConsulClusterId: d.Get("consul_cluster_id").(string),
        FirewallWhitelistIpv4: d.Get("firewall_whitelist_ipv4").(string),
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/cockroachdb/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudCockroachdbResponse McloudCockroachdbResponse
	err = json.Unmarshal(body, &mcloudCockroachdbResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("access_key_primary", mcloudCockroachdbResponse.AccessKeyPrimary)
    d.Set("access_key_secondary", mcloudCockroachdbResponse.AccessKeySecondary)
    d.Set("consul_cluster_id", mcloudCockroachdbResponse.ConsulClusterId)
    d.Set("firewall_whitelist_ipv4", mcloudCockroachdbResponse.FirewallWhitelistIpv4)
    d.Set("master_domain", mcloudCockroachdbResponse.MasterDomain)
    d.Set("name", mcloudCockroachdbResponse.Name)
    d.Set("pki_ca_id", mcloudCockroachdbResponse.PkiCaId)
    d.Set("server_pool_id", mcloudCockroachdbResponse.ServerPoolId)
    d.Set("status", mcloudCockroachdbResponse.Status)
    d.Set("ui_basic_auth_password", mcloudCockroachdbResponse.UiBasicAuthPassword)
    d.Set("ui_basic_auth_user", mcloudCockroachdbResponse.UiBasicAuthUser)
    d.Set("version", mcloudCockroachdbResponse.Version)

	return diags
}

func resourceMcloudCockroachdbRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/cockroachdb/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_cockroachdb %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudCockroachdbResponse McloudCockroachdbResponse
	err = json.Unmarshal(body, &mcloudCockroachdbResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudCockroachdbResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("access_key_primary", mcloudCockroachdbResponse.AccessKeyPrimary)
    d.Set("access_key_secondary", mcloudCockroachdbResponse.AccessKeySecondary)
    d.Set("consul_cluster_id", mcloudCockroachdbResponse.ConsulClusterId)
    d.Set("firewall_whitelist_ipv4", mcloudCockroachdbResponse.FirewallWhitelistIpv4)
    d.Set("master_domain", mcloudCockroachdbResponse.MasterDomain)
    d.Set("name", mcloudCockroachdbResponse.Name)
    d.Set("pki_ca_id", mcloudCockroachdbResponse.PkiCaId)
    d.Set("server_pool_id", mcloudCockroachdbResponse.ServerPoolId)
    d.Set("status", mcloudCockroachdbResponse.Status)
    d.Set("ui_basic_auth_password", mcloudCockroachdbResponse.UiBasicAuthPassword)
    d.Set("ui_basic_auth_user", mcloudCockroachdbResponse.UiBasicAuthUser)
    d.Set("version", mcloudCockroachdbResponse.Version)

	return diags
}
func resourceMcloudCockroachdbUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudCockroachdbCreate(ctx, d, m)
}

func resourceMcloudCockroachdbDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/cockroachdb/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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