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

type McloudConsulLoadbalancer struct {
    IpScopeAdminId string `json:"ip_scope_admin_id"`
    Name string `json:"name"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
}

type McloudConsulLoadbalancerResponse struct {
    IpScopeAdminId string `json:"ip_scope_admin_id"`
    Name string `json:"name"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
}

func resourceMcloudConsulLoadbalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudConsulLoadbalancerCreate,
		ReadContext:   resourceMcloudConsulLoadbalancerRead,
		UpdateContext: resourceMcloudConsulLoadbalancerUpdate,
		DeleteContext: resourceMcloudConsulLoadbalancerDelete,
		Schema: map[string]*schema.Schema{
			"ip_scope_admin_id": &schema.Schema{
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
			"server_pool_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
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
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMcloudConsulLoadbalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudConsulLoadbalancer{
        IpScopeAdminId: d.Get("ip_scope_admin_id").(string),
        Name: d.Get("name").(string),
        ServerPoolId: d.Get("server_pool_id").(string),
        Status: d.Get("status").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/consul-loadbalancer/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudConsulLoadbalancerResponse McloudConsulLoadbalancerResponse
	err = json.Unmarshal(body, &mcloudConsulLoadbalancerResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("ip_scope_admin_id", mcloudConsulLoadbalancerResponse.IpScopeAdminId)
    d.Set("name", mcloudConsulLoadbalancerResponse.Name)
    d.Set("server_pool_id", mcloudConsulLoadbalancerResponse.ServerPoolId)
    d.Set("status", mcloudConsulLoadbalancerResponse.Status)

	return diags
}

func resourceMcloudConsulLoadbalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/consul-loadbalancer/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_consul_loadbalancer %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudConsulLoadbalancerResponse McloudConsulLoadbalancerResponse
	err = json.Unmarshal(body, &mcloudConsulLoadbalancerResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudConsulLoadbalancerResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("ip_scope_admin_id", mcloudConsulLoadbalancerResponse.IpScopeAdminId)
    d.Set("name", mcloudConsulLoadbalancerResponse.Name)
    d.Set("server_pool_id", mcloudConsulLoadbalancerResponse.ServerPoolId)
    d.Set("status", mcloudConsulLoadbalancerResponse.Status)

	return diags
}
func resourceMcloudConsulLoadbalancerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudConsulLoadbalancerCreate(ctx, d, m)
}

func resourceMcloudConsulLoadbalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/consul-loadbalancer/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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