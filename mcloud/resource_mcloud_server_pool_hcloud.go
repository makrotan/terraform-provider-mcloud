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

type McloudServerPoolHcloud struct {
    InstanceCount int `json:"instance_count"`
    InstanceType string `json:"instance_type"`
    IpBlockId string `json:"ip_block_id,omitempty"`
    Location string `json:"location"`
    Name string `json:"name"`
    Status string `json:"status"`
}

type McloudServerPoolHcloudResponse struct {
    InstanceCount int `json:"instance_count"`
    InstanceType string `json:"instance_type"`
    IpBlockId string `json:"ip_block_id"`
    Location string `json:"location"`
    Name string `json:"name"`
    Status string `json:"status"`
}

func resourceMcloudServerPoolHcloud() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudServerPoolHcloudCreate,
		ReadContext:   resourceMcloudServerPoolHcloudRead,
		UpdateContext: resourceMcloudServerPoolHcloudUpdate,
		DeleteContext: resourceMcloudServerPoolHcloudDelete,
		Schema: map[string]*schema.Schema{
			"instance_count": &schema.Schema{
			    Type:     schema.TypeInt,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"instance_type": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"ip_block_id": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"location": &schema.Schema{
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

func resourceMcloudServerPoolHcloudCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudServerPoolHcloud{
        InstanceCount: d.Get("instance_count").(int),
        InstanceType: d.Get("instance_type").(string),
        Location: d.Get("location").(string),
        Name: d.Get("name").(string),
        Status: d.Get("status").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/server-pool-hcloud/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudServerPoolHcloudResponse McloudServerPoolHcloudResponse
	err = json.Unmarshal(body, &mcloudServerPoolHcloudResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("instance_count", mcloudServerPoolHcloudResponse.InstanceCount)
    d.Set("instance_type", mcloudServerPoolHcloudResponse.InstanceType)
    d.Set("ip_block_id", mcloudServerPoolHcloudResponse.IpBlockId)
    d.Set("location", mcloudServerPoolHcloudResponse.Location)
    d.Set("name", mcloudServerPoolHcloudResponse.Name)
    d.Set("status", mcloudServerPoolHcloudResponse.Status)

	return diags
}

func resourceMcloudServerPoolHcloudRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/server-pool-hcloud/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_server_pool_hcloud %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudServerPoolHcloudResponse McloudServerPoolHcloudResponse
	err = json.Unmarshal(body, &mcloudServerPoolHcloudResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudServerPoolHcloudResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("instance_count", mcloudServerPoolHcloudResponse.InstanceCount)
    d.Set("instance_type", mcloudServerPoolHcloudResponse.InstanceType)
    d.Set("ip_block_id", mcloudServerPoolHcloudResponse.IpBlockId)
    d.Set("location", mcloudServerPoolHcloudResponse.Location)
    d.Set("name", mcloudServerPoolHcloudResponse.Name)
    d.Set("status", mcloudServerPoolHcloudResponse.Status)

	return diags
}
func resourceMcloudServerPoolHcloudUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudServerPoolHcloudCreate(ctx, d, m)
}

func resourceMcloudServerPoolHcloudDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/server-pool-hcloud/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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