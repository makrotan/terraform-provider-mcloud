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

type McloudGrafana struct {
    AdminPassword string `json:"admin_password,omitempty"`
    Fqdn string `json:"fqdn"`
    Name string `json:"name"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
    Version string `json:"version"`
}

type McloudGrafanaResponse struct {
    AdminPassword string `json:"admin_password"`
    Fqdn string `json:"fqdn"`
    Name string `json:"name"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
    Version string `json:"version"`
}

func resourceMcloudGrafana() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudGrafanaCreate,
		ReadContext:   resourceMcloudGrafanaRead,
		UpdateContext: resourceMcloudGrafanaUpdate,
		DeleteContext: resourceMcloudGrafanaDelete,
		Schema: map[string]*schema.Schema{
			"admin_password": &schema.Schema{
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

func resourceMcloudGrafanaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudGrafana{
        Fqdn: d.Get("fqdn").(string),
        Name: d.Get("name").(string),
        ServerPoolId: d.Get("server_pool_id").(string),
        Status: d.Get("status").(string),
        Version: d.Get("version").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/grafana/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudGrafanaResponse McloudGrafanaResponse
	err = json.Unmarshal(body, &mcloudGrafanaResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("admin_password", mcloudGrafanaResponse.AdminPassword)
    d.Set("fqdn", mcloudGrafanaResponse.Fqdn)
    d.Set("name", mcloudGrafanaResponse.Name)
    d.Set("server_pool_id", mcloudGrafanaResponse.ServerPoolId)
    d.Set("status", mcloudGrafanaResponse.Status)
    d.Set("version", mcloudGrafanaResponse.Version)

	return diags
}

func resourceMcloudGrafanaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/grafana/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_grafana %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudGrafanaResponse McloudGrafanaResponse
	err = json.Unmarshal(body, &mcloudGrafanaResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudGrafanaResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("admin_password", mcloudGrafanaResponse.AdminPassword)
    d.Set("fqdn", mcloudGrafanaResponse.Fqdn)
    d.Set("name", mcloudGrafanaResponse.Name)
    d.Set("server_pool_id", mcloudGrafanaResponse.ServerPoolId)
    d.Set("status", mcloudGrafanaResponse.Status)
    d.Set("version", mcloudGrafanaResponse.Version)

	return diags
}
func resourceMcloudGrafanaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudGrafanaCreate(ctx, d, m)
}

func resourceMcloudGrafanaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/grafana/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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