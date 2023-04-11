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

type McloudYugabytedbDatabase struct {
    Name string `json:"name"`
    Password string `json:"password,omitempty"`
    Status string `json:"status"`
    Username string `json:"username,omitempty"`
    YugabytedbId string `json:"yugabytedb_id"`
}

type McloudYugabytedbDatabaseResponse struct {
    Name string `json:"name"`
    Password string `json:"password"`
    Status string `json:"status"`
    Username string `json:"username"`
    YugabytedbId string `json:"yugabytedb_id"`
}

func resourceMcloudYugabytedbDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudYugabytedbDatabaseCreate,
		ReadContext:   resourceMcloudYugabytedbDatabaseRead,
		UpdateContext: resourceMcloudYugabytedbDatabaseUpdate,
		DeleteContext: resourceMcloudYugabytedbDatabaseDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"password": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"status": &schema.Schema{
                Type:     schema.TypeString,
                Default: "running",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"username": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"yugabytedb_id": &schema.Schema{
                Type:     schema.TypeString,
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

func resourceMcloudYugabytedbDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudYugabytedbDatabase{
        Name: d.Get("name").(string),
        Status: d.Get("status").(string),
        YugabytedbId: d.Get("yugabytedb_id").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/yugabytedb-database/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudYugabytedbDatabaseResponse McloudYugabytedbDatabaseResponse
	err = json.Unmarshal(body, &mcloudYugabytedbDatabaseResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("name", mcloudYugabytedbDatabaseResponse.Name)
    d.Set("password", mcloudYugabytedbDatabaseResponse.Password)
    d.Set("status", mcloudYugabytedbDatabaseResponse.Status)
    d.Set("username", mcloudYugabytedbDatabaseResponse.Username)
    d.Set("yugabytedb_id", mcloudYugabytedbDatabaseResponse.YugabytedbId)

	return diags
}

func resourceMcloudYugabytedbDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/yugabytedb-database/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_yugabytedb_database %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudYugabytedbDatabaseResponse McloudYugabytedbDatabaseResponse
	err = json.Unmarshal(body, &mcloudYugabytedbDatabaseResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudYugabytedbDatabaseResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("name", mcloudYugabytedbDatabaseResponse.Name)
    d.Set("password", mcloudYugabytedbDatabaseResponse.Password)
    d.Set("status", mcloudYugabytedbDatabaseResponse.Status)
    d.Set("username", mcloudYugabytedbDatabaseResponse.Username)
    d.Set("yugabytedb_id", mcloudYugabytedbDatabaseResponse.YugabytedbId)

	return diags
}
func resourceMcloudYugabytedbDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudYugabytedbDatabaseCreate(ctx, d, m)
}

func resourceMcloudYugabytedbDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/yugabytedb-database/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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