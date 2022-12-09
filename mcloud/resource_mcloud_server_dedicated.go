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

type McloudServerDedicated struct {
    Name string `json:"name"`
    PoolName string `json:"pool_name"`
    Ipv4 string `json:"ipv4"`
    Provider string `json:"provider"`
    ProviderRef string `json:"provider_ref"`
    Region string `json:"region"`
    Az string `json:"az"`
}

type McloudServerDedicatedResponse struct {
    Name string `json:"name"`
    PoolName string `json:"pool_name"`
    Ipv4 string `json:"ipv4"`
    Provider string `json:"provider"`
    ProviderRef string `json:"provider_ref"`
    Region string `json:"region"`
    Az string `json:"az"`
}

func resourceMcloudServerDedicated() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudServerDedicatedCreate,
		ReadContext:   resourceMcloudServerDedicatedRead,
		UpdateContext: resourceMcloudServerDedicatedUpdate,
		DeleteContext: resourceMcloudServerDedicatedDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"pool_name": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"ipv4": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"server_provider": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"server_provider_ref": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"region": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"az": &schema.Schema{
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

func resourceMcloudServerDedicatedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudServerDedicated{
        Name: d.Get("name").(string),
        PoolName: d.Get("pool_name").(string),
        Ipv4: d.Get("ipv4").(string),
        Provider: d.Get("server_provider").(string),
        ProviderRef: d.Get("server_provider_ref").(string),
        Region: d.Get("region").(string),
        Az: d.Get("az").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/server-dedicated/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudServerDedicatedResponse McloudServerDedicatedResponse
	err = json.Unmarshal(body, &mcloudServerDedicatedResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("name", mcloudServerDedicatedResponse.Name)
    d.Set("pool_name", mcloudServerDedicatedResponse.PoolName)
    d.Set("ipv4", mcloudServerDedicatedResponse.Ipv4)
    d.Set("server_provider", mcloudServerDedicatedResponse.Provider)
    d.Set("server_provider_ref", mcloudServerDedicatedResponse.ProviderRef)
    d.Set("region", mcloudServerDedicatedResponse.Region)
    d.Set("az", mcloudServerDedicatedResponse.Az)

	return diags
}

func resourceMcloudServerDedicatedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/server-dedicated/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_server_dedicated %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudServerDedicatedResponse McloudServerDedicatedResponse
	err = json.Unmarshal(body, &mcloudServerDedicatedResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudServerDedicatedResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("name", mcloudServerDedicatedResponse.Name)
    d.Set("pool_name", mcloudServerDedicatedResponse.PoolName)
    d.Set("ipv4", mcloudServerDedicatedResponse.Ipv4)
    d.Set("server_provider", mcloudServerDedicatedResponse.Provider)
    d.Set("server_provider_ref", mcloudServerDedicatedResponse.ProviderRef)
    d.Set("region", mcloudServerDedicatedResponse.Region)
    d.Set("az", mcloudServerDedicatedResponse.Az)

	return diags
}
func resourceMcloudServerDedicatedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudServerDedicatedCreate(ctx, d, m)
}

func resourceMcloudServerDedicatedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/server-dedicated/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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