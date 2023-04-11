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

type McloudIpScopeBlockAssignment struct {
    BlockId string `json:"block_id"`
    Name string `json:"name"`
    ScopeId string `json:"scope_id"`
    Status string `json:"status"`
}

type McloudIpScopeBlockAssignmentResponse struct {
    BlockId string `json:"block_id"`
    Name string `json:"name"`
    ScopeId string `json:"scope_id"`
    Status string `json:"status"`
}

func resourceMcloudIpScopeBlockAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudIpScopeBlockAssignmentCreate,
		ReadContext:   resourceMcloudIpScopeBlockAssignmentRead,
		UpdateContext: resourceMcloudIpScopeBlockAssignmentUpdate,
		DeleteContext: resourceMcloudIpScopeBlockAssignmentDelete,
		Schema: map[string]*schema.Schema{
			"block_id": &schema.Schema{
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
			"scope_id": &schema.Schema{
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
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMcloudIpScopeBlockAssignmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudIpScopeBlockAssignment{
        BlockId: d.Get("block_id").(string),
        Name: d.Get("name").(string),
        ScopeId: d.Get("scope_id").(string),
        Status: d.Get("status").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ip-scope-block-assignment/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudIpScopeBlockAssignmentResponse McloudIpScopeBlockAssignmentResponse
	err = json.Unmarshal(body, &mcloudIpScopeBlockAssignmentResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("block_id", mcloudIpScopeBlockAssignmentResponse.BlockId)
    d.Set("name", mcloudIpScopeBlockAssignmentResponse.Name)
    d.Set("scope_id", mcloudIpScopeBlockAssignmentResponse.ScopeId)
    d.Set("status", mcloudIpScopeBlockAssignmentResponse.Status)

	return diags
}

func resourceMcloudIpScopeBlockAssignmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/ip-scope-block-assignment/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_ip_scope_block_assignment %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudIpScopeBlockAssignmentResponse McloudIpScopeBlockAssignmentResponse
	err = json.Unmarshal(body, &mcloudIpScopeBlockAssignmentResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudIpScopeBlockAssignmentResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("block_id", mcloudIpScopeBlockAssignmentResponse.BlockId)
    d.Set("name", mcloudIpScopeBlockAssignmentResponse.Name)
    d.Set("scope_id", mcloudIpScopeBlockAssignmentResponse.ScopeId)
    d.Set("status", mcloudIpScopeBlockAssignmentResponse.Status)

	return diags
}
func resourceMcloudIpScopeBlockAssignmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudIpScopeBlockAssignmentCreate(ctx, d, m)
}

func resourceMcloudIpScopeBlockAssignmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/ip-scope-block-assignment/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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