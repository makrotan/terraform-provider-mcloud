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

type McloudBackupPolicy struct {
	KeepBackupCount int    `json:"keep_backup_count"`
	Name            string `json:"name"`
	Ref             string `json:"ref"`
	Schedule        string `json:"schedule"`
	SchedulerJobId  string `json:"scheduler_job_id,omitempty"`
	Status          string `json:"status"`
}

type McloudBackupPolicyResponse struct {
	KeepBackupCount int    `json:"keep_backup_count"`
	Name            string `json:"name"`
	Ref             string `json:"ref"`
	Schedule        string `json:"schedule"`
	SchedulerJobId  string `json:"scheduler_job_id"`
	Status          string `json:"status"`
}

func resourceMcloudBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudBackupPolicyCreate,
		ReadContext:   resourceMcloudBackupPolicyRead,
		UpdateContext: resourceMcloudBackupPolicyUpdate,
		DeleteContext: resourceMcloudBackupPolicyDelete,
		Schema: map[string]*schema.Schema{
			"keep_backup_count": &schema.Schema{
				Type:     schema.TypeInt,
				Default:  7,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"schedule": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"scheduler_job_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: false, Computed: true, Optional: false, ForceNew: false,
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

func resourceMcloudBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudBackupPolicy{
		KeepBackupCount: d.Get("keep_backup_count").(int),
		Name:            d.Get("name").(string),
		Ref:             d.Get("ref").(string),
		Schedule:        d.Get("schedule").(string),
		Status:          d.Get("status").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/backup-policy/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudBackupPolicyResponse McloudBackupPolicyResponse
	err = json.Unmarshal(body, &mcloudBackupPolicyResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
	d.Set("keep_backup_count", mcloudBackupPolicyResponse.KeepBackupCount)
	d.Set("name", mcloudBackupPolicyResponse.Name)
	d.Set("ref", mcloudBackupPolicyResponse.Ref)
	d.Set("schedule", mcloudBackupPolicyResponse.Schedule)
	d.Set("scheduler_job_id", mcloudBackupPolicyResponse.SchedulerJobId)
	d.Set("status", mcloudBackupPolicyResponse.Status)

	return diags
}

func resourceMcloudBackupPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/backup-policy/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_backup_policy %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudBackupPolicyResponse McloudBackupPolicyResponse
	err = json.Unmarshal(body, &mcloudBackupPolicyResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudBackupPolicyResponse)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("keep_backup_count", mcloudBackupPolicyResponse.KeepBackupCount)
	d.Set("name", mcloudBackupPolicyResponse.Name)
	d.Set("ref", mcloudBackupPolicyResponse.Ref)
	d.Set("schedule", mcloudBackupPolicyResponse.Schedule)
	d.Set("scheduler_job_id", mcloudBackupPolicyResponse.SchedulerJobId)
	d.Set("status", mcloudBackupPolicyResponse.Status)

	return diags
}
func resourceMcloudBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudBackupPolicyCreate(ctx, d, m)
}

func resourceMcloudBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/backup-policy/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
