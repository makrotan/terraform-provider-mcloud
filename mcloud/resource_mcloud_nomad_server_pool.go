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

type McloudNomadServerPool struct {
    Name string `json:"name"`
    NomadClusterId string `json:"nomad_cluster_id"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
}

type McloudNomadServerPoolResponse struct {
    Name string `json:"name"`
    NomadClusterId string `json:"nomad_cluster_id"`
    ServerPoolId string `json:"server_pool_id"`
    Status string `json:"status"`
}

func resourceMcloudNomadServerPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudNomadServerPoolCreate,
		ReadContext:   resourceMcloudNomadServerPoolRead,
		UpdateContext: resourceMcloudNomadServerPoolUpdate,
		DeleteContext: resourceMcloudNomadServerPoolDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"nomad_cluster_id": &schema.Schema{
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

func resourceMcloudNomadServerPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudNomadServerPool{
        Name: d.Get("name").(string),
        NomadClusterId: d.Get("nomad_cluster_id").(string),
        ServerPoolId: d.Get("server_pool_id").(string),
        Status: d.Get("status").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/nomad-server-pool/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudNomadServerPoolResponse McloudNomadServerPoolResponse
	err = json.Unmarshal(body, &mcloudNomadServerPoolResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("name", mcloudNomadServerPoolResponse.Name)
    d.Set("nomad_cluster_id", mcloudNomadServerPoolResponse.NomadClusterId)
    d.Set("server_pool_id", mcloudNomadServerPoolResponse.ServerPoolId)
    d.Set("status", mcloudNomadServerPoolResponse.Status)

	return diags
}

func resourceMcloudNomadServerPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/nomad-server-pool/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_nomad_server_pool %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudNomadServerPoolResponse McloudNomadServerPoolResponse
	err = json.Unmarshal(body, &mcloudNomadServerPoolResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudNomadServerPoolResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("name", mcloudNomadServerPoolResponse.Name)
    d.Set("nomad_cluster_id", mcloudNomadServerPoolResponse.NomadClusterId)
    d.Set("server_pool_id", mcloudNomadServerPoolResponse.ServerPoolId)
    d.Set("status", mcloudNomadServerPoolResponse.Status)

	return diags
}
func resourceMcloudNomadServerPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudNomadServerPoolCreate(ctx, d, m)
}

func resourceMcloudNomadServerPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/nomad-server-pool/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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