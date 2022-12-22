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

type McloudMattermost struct {
    Name string `json:"name"`
    Fqdn string `json:"fqdn"`
    Sku string `json:"sku"`
    Version string `json:"version"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status,omitempty"`
}

type McloudMattermostResponse struct {
    Name string `json:"name"`
    Fqdn string `json:"fqdn"`
    Sku string `json:"sku"`
    Version string `json:"version"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
}

func resourceMcloudMattermost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudMattermostCreate,
		ReadContext:   resourceMcloudMattermostRead,
		UpdateContext: resourceMcloudMattermostUpdate,
		DeleteContext: resourceMcloudMattermostDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"fqdn": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"sku": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
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
			"server_pool_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: true,
			},
			"status": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMcloudMattermostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudMattermost{
        Name: d.Get("name").(string),
        Fqdn: d.Get("fqdn").(string),
        Sku: d.Get("sku").(string),
        Version: d.Get("version").(string),
        ServerPoolId: d.Get("server_pool_id").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/mattermost/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudMattermostResponse McloudMattermostResponse
	err = json.Unmarshal(body, &mcloudMattermostResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("name", mcloudMattermostResponse.Name)
    d.Set("fqdn", mcloudMattermostResponse.Fqdn)
    d.Set("sku", mcloudMattermostResponse.Sku)
    d.Set("version", mcloudMattermostResponse.Version)
    d.Set("server_pool_id", mcloudMattermostResponse.ServerPoolId)
    d.Set("status", mcloudMattermostResponse.Status)

	return diags
}

func resourceMcloudMattermostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/mattermost/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_mattermost %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudMattermostResponse McloudMattermostResponse
	err = json.Unmarshal(body, &mcloudMattermostResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudMattermostResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("name", mcloudMattermostResponse.Name)
    d.Set("fqdn", mcloudMattermostResponse.Fqdn)
    d.Set("sku", mcloudMattermostResponse.Sku)
    d.Set("version", mcloudMattermostResponse.Version)
    d.Set("server_pool_id", mcloudMattermostResponse.ServerPoolId)
    d.Set("status", mcloudMattermostResponse.Status)

	return diags
}
func resourceMcloudMattermostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudMattermostCreate(ctx, d, m)
}

func resourceMcloudMattermostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/mattermost/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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