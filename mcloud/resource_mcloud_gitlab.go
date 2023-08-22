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

type McloudGitlab struct {
	AdminPassword                  string                 `json:"admin_password,omitempty"`
	AdminUsername                  string                 `json:"admin_username,omitempty"`
	Fqdn                           string                 `json:"fqdn"`
	Meta                           map[string]interface{} `json:"meta"`
	Name                           string                 `json:"name"`
	ServerPoolId                   string                 `json:"server_pool_id"`
	SharedRunnersRegistrationToken string                 `json:"shared_runners_registration_token,omitempty"`
	Status                         string                 `json:"status"`
	Version                        string                 `json:"version"`
}

type McloudGitlabResponse struct {
	AdminPassword                  string                 `json:"admin_password"`
	AdminUsername                  string                 `json:"admin_username"`
	Fqdn                           string                 `json:"fqdn"`
	Meta                           map[string]interface{} `json:"meta"`
	Name                           string                 `json:"name"`
	ServerPoolId                   string                 `json:"server_pool_id"`
	SharedRunnersRegistrationToken string                 `json:"shared_runners_registration_token"`
	Status                         string                 `json:"status"`
	Version                        string                 `json:"version"`
}

func resourceMcloudGitlab() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudGitlabCreate,
		ReadContext:   resourceMcloudGitlabRead,
		UpdateContext: resourceMcloudGitlabUpdate,
		DeleteContext: resourceMcloudGitlabDelete,
		Schema: map[string]*schema.Schema{
			"admin_password": &schema.Schema{
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  false, Computed: true, Optional: false, ForceNew: false,
			},
			"admin_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"fqdn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"meta": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"shared_runners_registration_token": &schema.Schema{
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  false, Computed: true, Optional: false, ForceNew: false,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "running",
				Optional: true,
				Required: false,
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
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMcloudGitlabCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudGitlab{
		Fqdn:         d.Get("fqdn").(string),
		Meta:         d.Get("meta").(map[string]interface{}),
		Name:         d.Get("name").(string),
		ServerPoolId: d.Get("server_pool_id").(string),
		Status:       d.Get("status").(string),
		Version:      d.Get("version").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/gitlab/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudGitlabResponse McloudGitlabResponse
	err = json.Unmarshal(body, &mcloudGitlabResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
	d.Set("admin_password", mcloudGitlabResponse.AdminPassword)
	d.Set("admin_username", mcloudGitlabResponse.AdminUsername)
	d.Set("fqdn", mcloudGitlabResponse.Fqdn)
	d.Set("meta", mcloudGitlabResponse.Meta)
	d.Set("name", mcloudGitlabResponse.Name)
	d.Set("server_pool_id", mcloudGitlabResponse.ServerPoolId)
	d.Set("shared_runners_registration_token", mcloudGitlabResponse.SharedRunnersRegistrationToken)
	d.Set("status", mcloudGitlabResponse.Status)
	d.Set("version", mcloudGitlabResponse.Version)

	return diags
}

func resourceMcloudGitlabRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/gitlab/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_gitlab %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudGitlabResponse McloudGitlabResponse
	err = json.Unmarshal(body, &mcloudGitlabResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudGitlabResponse)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("admin_password", mcloudGitlabResponse.AdminPassword)
	d.Set("admin_username", mcloudGitlabResponse.AdminUsername)
	d.Set("fqdn", mcloudGitlabResponse.Fqdn)
	d.Set("meta", mcloudGitlabResponse.Meta)
	d.Set("name", mcloudGitlabResponse.Name)
	d.Set("server_pool_id", mcloudGitlabResponse.ServerPoolId)
	d.Set("shared_runners_registration_token", mcloudGitlabResponse.SharedRunnersRegistrationToken)
	d.Set("status", mcloudGitlabResponse.Status)
	d.Set("version", mcloudGitlabResponse.Version)

	return diags
}
func resourceMcloudGitlabUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudGitlabCreate(ctx, d, m)
}

func resourceMcloudGitlabDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/gitlab/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
