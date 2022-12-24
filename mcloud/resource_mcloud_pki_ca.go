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

type McloudPkiCa struct {
	Name         string `json:"name"`
	ValidDays    int    `json:"valid_days"`
	Country      string `json:"country"`
	State        string `json:"state"`
	City         string `json:"city"`
	Organisation string `json:"organisation"`
	Unit         string `json:"unit"`
	Email        string `json:"email"`
	KeyPub       string `json:"key_pub,omitempty"`
	KeyPriv      string `json:"key_priv,omitempty"`
}

type McloudPkiCaResponse struct {
	Name         string `json:"name"`
	ValidDays    int    `json:"valid_days"`
	Country      string `json:"country"`
	State        string `json:"state"`
	City         string `json:"city"`
	Organisation string `json:"organisation"`
	Unit         string `json:"unit"`
	Email        string `json:"email"`
	KeyPub       string `json:"key_pub"`
	KeyPriv      string `json:"key_priv"`
}

func resourceMcloudPkiCa() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudPkiCaCreate,
		ReadContext:   resourceMcloudPkiCaRead,
		UpdateContext: resourceMcloudPkiCaUpdate,
		DeleteContext: resourceMcloudPkiCaDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"valid_days": &schema.Schema{
				Type:     schema.TypeInt,
				Default:  7300,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"country": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "DE",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "DE",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"city": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "Leipzig",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"organisation": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "Makrotan",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"unit": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "IT",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "info@makrotan.com",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"key_pub": &schema.Schema{
				Type:     schema.TypeString,
				Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"key_priv": &schema.Schema{
				Type:     schema.TypeString,
				Required: false, Computed: true, Optional: false, ForceNew: false,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMcloudPkiCaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudPkiCa{
		Name:         d.Get("name").(string),
		ValidDays:    d.Get("valid_days").(int),
		Country:      d.Get("country").(string),
		State:        d.Get("state").(string),
		City:         d.Get("city").(string),
		Organisation: d.Get("organisation").(string),
		Unit:         d.Get("unit").(string),
		Email:        d.Get("email").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/pki-ca/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudPkiCaResponse McloudPkiCaResponse
	err = json.Unmarshal(body, &mcloudPkiCaResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
	d.Set("name", mcloudPkiCaResponse.Name)
	d.Set("valid_days", mcloudPkiCaResponse.ValidDays)
	d.Set("country", mcloudPkiCaResponse.Country)
	d.Set("state", mcloudPkiCaResponse.State)
	d.Set("city", mcloudPkiCaResponse.City)
	d.Set("organisation", mcloudPkiCaResponse.Organisation)
	d.Set("unit", mcloudPkiCaResponse.Unit)
	d.Set("email", mcloudPkiCaResponse.Email)
	d.Set("key_pub", mcloudPkiCaResponse.KeyPub)
	d.Set("key_priv", mcloudPkiCaResponse.KeyPriv)

	return diags
}

func resourceMcloudPkiCaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/pki-ca/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_pki_ca %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudPkiCaResponse McloudPkiCaResponse
	err = json.Unmarshal(body, &mcloudPkiCaResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudPkiCaResponse)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", mcloudPkiCaResponse.Name)
	d.Set("valid_days", mcloudPkiCaResponse.ValidDays)
	d.Set("country", mcloudPkiCaResponse.Country)
	d.Set("state", mcloudPkiCaResponse.State)
	d.Set("city", mcloudPkiCaResponse.City)
	d.Set("organisation", mcloudPkiCaResponse.Organisation)
	d.Set("unit", mcloudPkiCaResponse.Unit)
	d.Set("email", mcloudPkiCaResponse.Email)
	d.Set("key_pub", mcloudPkiCaResponse.KeyPub)
	d.Set("key_priv", mcloudPkiCaResponse.KeyPriv)

	return diags
}
func resourceMcloudPkiCaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudPkiCaCreate(ctx, d, m)
}

func resourceMcloudPkiCaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/pki-ca/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
