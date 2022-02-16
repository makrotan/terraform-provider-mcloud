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

func resourceK3sCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceK3sClusterCreate,
		ReadContext:   resourceK3sClusterRead,
		UpdateContext: resourceK3sClusterUpdate,
		DeleteContext: resourceK3sClusterDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_key_primary": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"sku": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"master_server_pool_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"k3s_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"firewall_whitelist_ipv4": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceK3sClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	cluster := K3SClusterRequest{
		Name:                  name,
		MasterServerPoolID:    d.Get("master_server_pool_id").(string),
		SKU:                   d.Get("sku").(string),
		K3sVersion:            d.Get("k3s_version").(string),
		FirewallWhitelistIPv4: d.Get("firewall_whitelist_ipv4").(string),
		RunSetup:              true,
	}

	rb, err := json.Marshal(cluster)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/k3s-cluster/%s", strings.Trim(c.HostURL, "/"), name), strings.NewReader(string(rb)))
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", c.Token)

	res, err := c.HTTPClient.Do(req)
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

	var k3sClusterResponse K3sClusterResponse
	err = json.Unmarshal(body, &k3sClusterResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	if !k3sClusterResponse.Success {
		return diag.FromErr(fmt.Errorf(k3sClusterResponse.Error))
	}

	//debug(strconv.Itoa(k3sClusterResponse.Task.Id))
	err = c.waitForTaskToFinish(k3sClusterResponse.Task.Id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	d.Set("master_server_pool_id", k3sClusterResponse.K3SCluster.MasterServerPoolID)
	d.Set("sku", k3sClusterResponse.K3SCluster.SKU)
	d.Set("k3s_version", k3sClusterResponse.K3SCluster.K3sVersion)
	d.Set("firewall_whitelist_ipv4", k3sClusterResponse.K3SCluster.FirewallWhitelistIPv4)
	d.Set("access_key_primary", k3sClusterResponse.K3SCluster.AccessKeyPrimary)

	return diags
}

func resourceK3sClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	client := c.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/k3s-cluster/%s", strings.Trim(c.HostURL, "/"), name), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", c.Token)

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
		log.Printf("[WARN] k3s cluster %s does not exist", name)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var k3sClusterResponse K3sClusterResponse
	err = json.Unmarshal(body, &k3sClusterResponse)
	//err = json.NewDecoder(resp.Body).Decode(sshKeyResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	if !k3sClusterResponse.Success {
		return diag.FromErr(fmt.Errorf(k3sClusterResponse.Error))
	}

	d.Set("master_server_pool_id", k3sClusterResponse.K3SCluster.MasterServerPoolID)
	d.Set("sku", k3sClusterResponse.K3SCluster.SKU)
	d.Set("k3s_version", k3sClusterResponse.K3SCluster.K3sVersion)
	d.Set("firewall_whitelist_ipv4", k3sClusterResponse.K3SCluster.FirewallWhitelistIPv4)
	d.Set("access_key_primary", k3sClusterResponse.K3SCluster.AccessKeyPrimary)

	return diags
}

func resourceK3sClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceK3sClusterCreate(ctx, d, m)
}

func resourceK3sClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	client := c.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/k3s-cluster/%s", strings.Trim(c.HostURL, "/"), name), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", c.Token)

	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	return diags
}
