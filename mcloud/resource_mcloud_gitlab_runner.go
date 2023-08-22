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

type McloudGitlabRunner struct {
	GitlabId     string                 `json:"gitlab_id"`
	Meta         map[string]interface{} `json:"meta"`
	Name         string                 `json:"name"`
	ServerPoolId string                 `json:"server_pool_id"`
	Status       string                 `json:"status"`
	Tags         string                 `json:"tags"`
	Version      string                 `json:"version"`
}

type McloudGitlabRunnerResponse struct {
	GitlabId     string                 `json:"gitlab_id"`
	Meta         map[string]interface{} `json:"meta"`
	Name         string                 `json:"name"`
	ServerPoolId string                 `json:"server_pool_id"`
	Status       string                 `json:"status"`
	Tags         string                 `json:"tags"`
	Version      string                 `json:"version"`
}

func resourceMcloudGitlabRunner() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudGitlabRunnerCreate,
		ReadContext:   resourceMcloudGitlabRunnerRead,
		UpdateContext: resourceMcloudGitlabRunnerUpdate,
		DeleteContext: resourceMcloudGitlabRunnerDelete,
		Schema: map[string]*schema.Schema{
			"gitlab_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
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
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "running",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeString,
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

func resourceMcloudGitlabRunnerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudGitlabRunner{
		GitlabId:     d.Get("gitlab_id").(string),
		Meta:         d.Get("meta").(map[string]interface{}),
		Name:         d.Get("name").(string),
		ServerPoolId: d.Get("server_pool_id").(string),
		Status:       d.Get("status").(string),
		Tags:         d.Get("tags").(string),
		Version:      d.Get("version").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/gitlab-runner/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudGitlabRunnerResponse McloudGitlabRunnerResponse
	err = json.Unmarshal(body, &mcloudGitlabRunnerResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
	d.Set("gitlab_id", mcloudGitlabRunnerResponse.GitlabId)
	d.Set("meta", mcloudGitlabRunnerResponse.Meta)
	d.Set("name", mcloudGitlabRunnerResponse.Name)
	d.Set("server_pool_id", mcloudGitlabRunnerResponse.ServerPoolId)
	d.Set("status", mcloudGitlabRunnerResponse.Status)
	d.Set("tags", mcloudGitlabRunnerResponse.Tags)
	d.Set("version", mcloudGitlabRunnerResponse.Version)

	return diags
}

func resourceMcloudGitlabRunnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/gitlab-runner/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_gitlab_runner %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudGitlabRunnerResponse McloudGitlabRunnerResponse
	err = json.Unmarshal(body, &mcloudGitlabRunnerResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudGitlabRunnerResponse)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("gitlab_id", mcloudGitlabRunnerResponse.GitlabId)
	d.Set("meta", mcloudGitlabRunnerResponse.Meta)
	d.Set("name", mcloudGitlabRunnerResponse.Name)
	d.Set("server_pool_id", mcloudGitlabRunnerResponse.ServerPoolId)
	d.Set("status", mcloudGitlabRunnerResponse.Status)
	d.Set("tags", mcloudGitlabRunnerResponse.Tags)
	d.Set("version", mcloudGitlabRunnerResponse.Version)

	return diags
}
func resourceMcloudGitlabRunnerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudGitlabRunnerCreate(ctx, d, m)
}

func resourceMcloudGitlabRunnerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/gitlab-runner/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
