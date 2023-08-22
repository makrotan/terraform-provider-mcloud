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
	Az            string `json:"az"`
	CpuCores      int    `json:"cpu_cores"`
	DiskSize      int    `json:"disk_size"`
	Ipv4          string `json:"ipv4"`
	Ipv6          string `json:"ipv6"`
	Memory        int    `json:"memory"`
	Name          string `json:"name"`
	PoolId        string `json:"pool_id"`
	PricePerMonth int    `json:"price_per_month"`
	Provider      string `json:"provider"`
	ProviderRef   string `json:"provider_ref"`
	Region        string `json:"region"`
	Status        string `json:"status"`
}

type McloudServerDedicatedResponse struct {
	Az            string `json:"az"`
	CpuCores      int    `json:"cpu_cores"`
	DiskSize      int    `json:"disk_size"`
	Ipv4          string `json:"ipv4"`
	Ipv6          string `json:"ipv6"`
	Memory        int    `json:"memory"`
	Name          string `json:"name"`
	PoolId        string `json:"pool_id"`
	PricePerMonth int    `json:"price_per_month"`
	Provider      string `json:"provider"`
	ProviderRef   string `json:"provider_ref"`
	Region        string `json:"region"`
	Status        string `json:"status"`
}

func resourceMcloudServerDedicated() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudServerDedicatedCreate,
		ReadContext:   resourceMcloudServerDedicatedRead,
		UpdateContext: resourceMcloudServerDedicatedUpdate,
		DeleteContext: resourceMcloudServerDedicatedDelete,
		Schema: map[string]*schema.Schema{
			"az": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"cpu_cores": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"disk_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
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
			"ipv6": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"pool_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"price_per_month": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"provider_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"provider_ref": &schema.Schema{
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
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "running",
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

func resourceMcloudServerDedicatedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudServerDedicated{
		Az:            d.Get("az").(string),
		CpuCores:      d.Get("cpu_cores").(int),
		DiskSize:      d.Get("disk_size").(int),
		Ipv4:          d.Get("ipv4").(string),
		Ipv6:          d.Get("ipv6").(string),
		Memory:        d.Get("memory").(int),
		Name:          d.Get("name").(string),
		PoolId:        d.Get("pool_id").(string),
		PricePerMonth: d.Get("price_per_month").(int),
		Provider:      d.Get("provider_id").(string),
		ProviderRef:   d.Get("provider_ref").(string),
		Region:        d.Get("region").(string),
		Status:        d.Get("status").(string),
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
	d.Set("az", mcloudServerDedicatedResponse.Az)
	d.Set("cpu_cores", mcloudServerDedicatedResponse.CpuCores)
	d.Set("disk_size", mcloudServerDedicatedResponse.DiskSize)
	d.Set("ipv4", mcloudServerDedicatedResponse.Ipv4)
	d.Set("ipv6", mcloudServerDedicatedResponse.Ipv6)
	d.Set("memory", mcloudServerDedicatedResponse.Memory)
	d.Set("name", mcloudServerDedicatedResponse.Name)
	d.Set("pool_id", mcloudServerDedicatedResponse.PoolId)
	d.Set("price_per_month", mcloudServerDedicatedResponse.PricePerMonth)
	d.Set("provider_id", mcloudServerDedicatedResponse.Provider)
	d.Set("provider_ref", mcloudServerDedicatedResponse.ProviderRef)
	d.Set("region", mcloudServerDedicatedResponse.Region)
	d.Set("status", mcloudServerDedicatedResponse.Status)

	return diags
}

func resourceMcloudServerDedicatedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/server-dedicated/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
	d.Set("az", mcloudServerDedicatedResponse.Az)
	d.Set("cpu_cores", mcloudServerDedicatedResponse.CpuCores)
	d.Set("disk_size", mcloudServerDedicatedResponse.DiskSize)
	d.Set("ipv4", mcloudServerDedicatedResponse.Ipv4)
	d.Set("ipv6", mcloudServerDedicatedResponse.Ipv6)
	d.Set("memory", mcloudServerDedicatedResponse.Memory)
	d.Set("name", mcloudServerDedicatedResponse.Name)
	d.Set("pool_id", mcloudServerDedicatedResponse.PoolId)
	d.Set("price_per_month", mcloudServerDedicatedResponse.PricePerMonth)
	d.Set("provider_id", mcloudServerDedicatedResponse.Provider)
	d.Set("provider_ref", mcloudServerDedicatedResponse.ProviderRef)
	d.Set("region", mcloudServerDedicatedResponse.Region)
	d.Set("status", mcloudServerDedicatedResponse.Status)

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
