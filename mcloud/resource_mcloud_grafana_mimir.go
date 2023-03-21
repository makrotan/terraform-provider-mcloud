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

type McloudGrafanaMimir struct {
    AccessKey string `json:"access_key,omitempty"`
    BasicAuthPassword string `json:"basic_auth_password,omitempty"`
    BasicAuthUser string `json:"basic_auth_user,omitempty"`
    Fqdn string `json:"fqdn"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
    Version string `json:"version"`
}

type McloudGrafanaMimirResponse struct {
    AccessKey string `json:"access_key"`
    BasicAuthPassword string `json:"basic_auth_password"`
    BasicAuthUser string `json:"basic_auth_user"`
    Fqdn string `json:"fqdn"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
    Version string `json:"version"`
}

func resourceMcloudGrafanaMimir() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudGrafanaMimirCreate,
		ReadContext:   resourceMcloudGrafanaMimirRead,
		UpdateContext: resourceMcloudGrafanaMimirUpdate,
		DeleteContext: resourceMcloudGrafanaMimirDelete,
		Schema: map[string]*schema.Schema{
			"access_key": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"basic_auth_password": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"basic_auth_user": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"fqdn": &schema.Schema{
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
				ForceNew: true,
			},
			"server_pool_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
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
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMcloudGrafanaMimirCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudGrafanaMimir{
        Fqdn: d.Get("fqdn").(string),
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/grafana-mimir/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudGrafanaMimirResponse McloudGrafanaMimirResponse
	err = json.Unmarshal(body, &mcloudGrafanaMimirResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("access_key", mcloudGrafanaMimirResponse.AccessKey)
    d.Set("basic_auth_password", mcloudGrafanaMimirResponse.BasicAuthPassword)
    d.Set("basic_auth_user", mcloudGrafanaMimirResponse.BasicAuthUser)
    d.Set("fqdn", mcloudGrafanaMimirResponse.Fqdn)
    d.Set("name", mcloudGrafanaMimirResponse.Name)
    d.Set("pki_ca_id", mcloudGrafanaMimirResponse.PkiCaId)
    d.Set("server_pool_id", mcloudGrafanaMimirResponse.ServerPoolId)
    d.Set("status", mcloudGrafanaMimirResponse.Status)
    d.Set("version", mcloudGrafanaMimirResponse.Version)

	return diags
}

func resourceMcloudGrafanaMimirRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/grafana-mimir/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_grafana_mimir %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudGrafanaMimirResponse McloudGrafanaMimirResponse
	err = json.Unmarshal(body, &mcloudGrafanaMimirResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudGrafanaMimirResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("access_key", mcloudGrafanaMimirResponse.AccessKey)
    d.Set("basic_auth_password", mcloudGrafanaMimirResponse.BasicAuthPassword)
    d.Set("basic_auth_user", mcloudGrafanaMimirResponse.BasicAuthUser)
    d.Set("fqdn", mcloudGrafanaMimirResponse.Fqdn)
    d.Set("name", mcloudGrafanaMimirResponse.Name)
    d.Set("pki_ca_id", mcloudGrafanaMimirResponse.PkiCaId)
    d.Set("server_pool_id", mcloudGrafanaMimirResponse.ServerPoolId)
    d.Set("status", mcloudGrafanaMimirResponse.Status)
    d.Set("version", mcloudGrafanaMimirResponse.Version)

	return diags
}
func resourceMcloudGrafanaMimirUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudGrafanaMimirCreate(ctx, d, m)
}

func resourceMcloudGrafanaMimirDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/grafana-mimir/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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