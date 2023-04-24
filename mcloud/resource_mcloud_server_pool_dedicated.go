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

type McloudServerPoolDedicated struct {
    ConsulClusterId string `json:"consul_cluster_id"`
    IpBlockId string `json:"ip_block_id,omitempty"`
    Name string `json:"name"`
    Status string `json:"status"`
}

type McloudServerPoolDedicatedResponse struct {
    ConsulClusterId string `json:"consul_cluster_id"`
    IpBlockId string `json:"ip_block_id"`
    Name string `json:"name"`
    Status string `json:"status"`
}

func resourceMcloudServerPoolDedicated() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudServerPoolDedicatedCreate,
		ReadContext:   resourceMcloudServerPoolDedicatedRead,
		UpdateContext: resourceMcloudServerPoolDedicatedUpdate,
		DeleteContext: resourceMcloudServerPoolDedicatedDelete,
		Schema: map[string]*schema.Schema{
			"consul_cluster_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"ip_block_id": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
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

func resourceMcloudServerPoolDedicatedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudServerPoolDedicated{
        ConsulClusterId: d.Get("consul_cluster_id").(string),
        Name: d.Get("name").(string),
        Status: d.Get("status").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/server-pool-dedicated/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudServerPoolDedicatedResponse McloudServerPoolDedicatedResponse
	err = json.Unmarshal(body, &mcloudServerPoolDedicatedResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("consul_cluster_id", mcloudServerPoolDedicatedResponse.ConsulClusterId)
    d.Set("ip_block_id", mcloudServerPoolDedicatedResponse.IpBlockId)
    d.Set("name", mcloudServerPoolDedicatedResponse.Name)
    d.Set("status", mcloudServerPoolDedicatedResponse.Status)

	return diags
}

func resourceMcloudServerPoolDedicatedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/server-pool-dedicated/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_server_pool_dedicated %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudServerPoolDedicatedResponse McloudServerPoolDedicatedResponse
	err = json.Unmarshal(body, &mcloudServerPoolDedicatedResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudServerPoolDedicatedResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("consul_cluster_id", mcloudServerPoolDedicatedResponse.ConsulClusterId)
    d.Set("ip_block_id", mcloudServerPoolDedicatedResponse.IpBlockId)
    d.Set("name", mcloudServerPoolDedicatedResponse.Name)
    d.Set("status", mcloudServerPoolDedicatedResponse.Status)

	return diags
}
func resourceMcloudServerPoolDedicatedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudServerPoolDedicatedCreate(ctx, d, m)
}

func resourceMcloudServerPoolDedicatedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/server-pool-dedicated/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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