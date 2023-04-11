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
    AccessKeyPrimary string `json:"access_key_primary,omitempty"`
    AccessKeySecondary string `json:"access_key_secondary,omitempty"`
    IpScopeId string `json:"ip_scope_id"`
    MasterServerPoolId string `json:"master_server_pool_id"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    RootToken string `json:"root_token"`
    SealKeys string `json:"seal_keys"`
    Status string `json:"status"`
    Version string `json:"version"`
}

type McloudVaultClusterResponse struct {
    AccessKeyPrimary string `json:"access_key_primary"`
    AccessKeySecondary string `json:"access_key_secondary"`
    IpScopeId string `json:"ip_scope_id"`
    MasterServerPoolId string `json:"master_server_pool_id"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    RootToken string `json:"root_token"`
    SealKeys string `json:"seal_keys"`
    Status string `json:"status"`
    Version string `json:"version"`
}

func resourceMcloudVaultCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudVaultClusterCreate,
		ReadContext:   resourceMcloudVaultClusterRead,
		UpdateContext: resourceMcloudVaultClusterUpdate,
		DeleteContext: resourceMcloudVaultClusterDelete,
		Schema: map[string]*schema.Schema{
			"access_key_primary": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"access_key_secondary": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"ip_scope_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"master_server_pool_id": &schema.Schema{
                Type:     schema.TypeString,
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
			"root_token": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"seal_keys": &schema.Schema{
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

func resourceMcloudVaultClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudVaultCluster{
        IpScopeId: d.Get("ip_scope_id").(string),
        MasterServerPoolId: d.Get("master_server_pool_id").(string),
        Name: d.Get("name").(string),
        PkiCaId: d.Get("pki_ca_id").(string),
        RootToken: d.Get("root_token").(string),
        SealKeys: d.Get("seal_keys").(string),
        Status: d.Get("status").(string),
        Version: d.Get("version").(string),
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
    d.Set("access_key_primary", mcloudVaultClusterResponse.AccessKeyPrimary)
    d.Set("access_key_secondary", mcloudVaultClusterResponse.AccessKeySecondary)
    d.Set("ip_scope_id", mcloudVaultClusterResponse.IpScopeId)
    d.Set("master_server_pool_id", mcloudVaultClusterResponse.MasterServerPoolId)
    d.Set("name", mcloudVaultClusterResponse.Name)
    d.Set("pki_ca_id", mcloudVaultClusterResponse.PkiCaId)
    d.Set("root_token", mcloudVaultClusterResponse.RootToken)
    d.Set("seal_keys", mcloudVaultClusterResponse.SealKeys)
    d.Set("status", mcloudVaultClusterResponse.Status)
    d.Set("version", mcloudVaultClusterResponse.Version)

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
    d.Set("access_key_primary", mcloudVaultClusterResponse.AccessKeyPrimary)
    d.Set("access_key_secondary", mcloudVaultClusterResponse.AccessKeySecondary)
    d.Set("ip_scope_id", mcloudVaultClusterResponse.IpScopeId)
    d.Set("master_server_pool_id", mcloudVaultClusterResponse.MasterServerPoolId)
    d.Set("name", mcloudVaultClusterResponse.Name)
    d.Set("pki_ca_id", mcloudVaultClusterResponse.PkiCaId)
    d.Set("root_token", mcloudVaultClusterResponse.RootToken)
    d.Set("seal_keys", mcloudVaultClusterResponse.SealKeys)
    d.Set("status", mcloudVaultClusterResponse.Status)
    d.Set("version", mcloudVaultClusterResponse.Version)

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