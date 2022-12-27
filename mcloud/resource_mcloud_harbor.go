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

type McloudHarbor struct {
	Name          string `json:"name"`
	Fqdn          string `json:"fqdn"`
	Sku           string `json:"sku"`
	Version       string `json:"version"`
	ServerPoolId  string `json:"server_pool_id"`
	Status        string `json:"status"`
	AdminPassword string `json:"admin_password,omitempty"`
}

type McloudHarborResponse struct {
	Name          string `json:"name"`
	Fqdn          string `json:"fqdn"`
	Sku           string `json:"sku"`
	Version       string `json:"version"`
	ServerPoolId  string `json:"server_pool_id"`
	Status        string `json:"status"`
	AdminPassword string `json:"admin_password"`
}

func resourceMcloudHarbor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudHarborCreate,
		ReadContext:   resourceMcloudHarborRead,
		UpdateContext: resourceMcloudHarborUpdate,
		DeleteContext: resourceMcloudHarborDelete,
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
				Default:  "running",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"admin_password": &schema.Schema{
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  false, Computed: true, Optional: false, ForceNew: false,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMcloudHarborCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudHarbor{
		Name:         d.Get("name").(string),
		Fqdn:         d.Get("fqdn").(string),
		Sku:          d.Get("sku").(string),
		Version:      d.Get("version").(string),
		ServerPoolId: d.Get("server_pool_id").(string),
		Status:       d.Get("status").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/harbor/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudHarborResponse McloudHarborResponse
	err = json.Unmarshal(body, &mcloudHarborResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
	d.Set("name", mcloudHarborResponse.Name)
	d.Set("fqdn", mcloudHarborResponse.Fqdn)
	d.Set("sku", mcloudHarborResponse.Sku)
	d.Set("version", mcloudHarborResponse.Version)
	d.Set("server_pool_id", mcloudHarborResponse.ServerPoolId)
	d.Set("status", mcloudHarborResponse.Status)
	d.Set("admin_password", mcloudHarborResponse.AdminPassword)

	return diags
}

func resourceMcloudHarborRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/harbor/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_harbor %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudHarborResponse McloudHarborResponse
	err = json.Unmarshal(body, &mcloudHarborResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudHarborResponse)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", mcloudHarborResponse.Name)
	d.Set("fqdn", mcloudHarborResponse.Fqdn)
	d.Set("sku", mcloudHarborResponse.Sku)
	d.Set("version", mcloudHarborResponse.Version)
	d.Set("server_pool_id", mcloudHarborResponse.ServerPoolId)
	d.Set("status", mcloudHarborResponse.Status)
	d.Set("admin_password", mcloudHarborResponse.AdminPassword)

	return diags
}
func resourceMcloudHarborUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudHarborCreate(ctx, d, m)
}

func resourceMcloudHarborDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/harbor/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
