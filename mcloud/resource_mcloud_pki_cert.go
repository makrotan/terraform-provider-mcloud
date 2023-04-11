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

type McloudPkiCert struct {
    CaId string `json:"ca_id"`
    CommonName string `json:"common_name"`
    KeyPriv string `json:"key_priv,omitempty"`
    KeyPub string `json:"key_pub,omitempty"`
    Name string `json:"name"`
}

type McloudPkiCertResponse struct {
    CaId string `json:"ca_id"`
    CommonName string `json:"common_name"`
    KeyPriv string `json:"key_priv"`
    KeyPub string `json:"key_pub"`
    Name string `json:"name"`
}

func resourceMcloudPkiCert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudPkiCertCreate,
		ReadContext:   resourceMcloudPkiCertRead,
		UpdateContext: resourceMcloudPkiCertUpdate,
		DeleteContext: resourceMcloudPkiCertDelete,
		Schema: map[string]*schema.Schema{
			"ca_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: true,
			},
			"common_name": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"key_priv": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"key_pub": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMcloudPkiCertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudPkiCert{
        CaId: d.Get("ca_id").(string),
        CommonName: d.Get("common_name").(string),
        Name: d.Get("name").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/pki-cert/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudPkiCertResponse McloudPkiCertResponse
	err = json.Unmarshal(body, &mcloudPkiCertResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("ca_id", mcloudPkiCertResponse.CaId)
    d.Set("common_name", mcloudPkiCertResponse.CommonName)
    d.Set("key_priv", mcloudPkiCertResponse.KeyPriv)
    d.Set("key_pub", mcloudPkiCertResponse.KeyPub)
    d.Set("name", mcloudPkiCertResponse.Name)

	return diags
}

func resourceMcloudPkiCertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/pki-cert/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_pki_cert %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudPkiCertResponse McloudPkiCertResponse
	err = json.Unmarshal(body, &mcloudPkiCertResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudPkiCertResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("ca_id", mcloudPkiCertResponse.CaId)
    d.Set("common_name", mcloudPkiCertResponse.CommonName)
    d.Set("key_priv", mcloudPkiCertResponse.KeyPriv)
    d.Set("key_pub", mcloudPkiCertResponse.KeyPub)
    d.Set("name", mcloudPkiCertResponse.Name)

	return diags
}
func resourceMcloudPkiCertUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudPkiCertCreate(ctx, d, m)
}

func resourceMcloudPkiCertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/pki-cert/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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