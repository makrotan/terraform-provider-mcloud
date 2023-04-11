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

type McloudKeycloak struct {
    AdminPassword string `json:"admin_password,omitempty"`
    Fqdn string `json:"fqdn"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    SecretKey string `json:"secret_key,omitempty"`
    ServerPoolId string `json:"server_pool_id"`
    Sku string `json:"sku"`
    Status string `json:"status"`
    Themes string `json:"themes"`
    Version string `json:"version"`
}

type McloudKeycloakResponse struct {
    AdminPassword string `json:"admin_password"`
    Fqdn string `json:"fqdn"`
    Name string `json:"name"`
    PkiCaId string `json:"pki_ca_id"`
    SecretKey string `json:"secret_key"`
    ServerPoolId string `json:"server_pool_id"`
    Sku string `json:"sku"`
    Status string `json:"status"`
    Themes string `json:"themes"`
    Version string `json:"version"`
}

func resourceMcloudKeycloak() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMcloudKeycloakCreate,
		ReadContext:   resourceMcloudKeycloakRead,
		UpdateContext: resourceMcloudKeycloakUpdate,
		DeleteContext: resourceMcloudKeycloakDelete,
		Schema: map[string]*schema.Schema{
			"admin_password": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"fqdn": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"pki_ca_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"secret_key": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
			"server_pool_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"sku": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
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
			"themes": &schema.Schema{
                Type:     schema.TypeString,
                Sensitive: true,
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

func resourceMcloudKeycloakCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("name").(string)
	instance := McloudKeycloak{
        Fqdn: d.Get("fqdn").(string),
        Name: d.Get("name").(string),
        PkiCaId: d.Get("pki_ca_id").(string),
        ServerPoolId: d.Get("server_pool_id").(string),
        Sku: d.Get("sku").(string),
        Status: d.Get("status").(string),
        Themes: d.Get("themes").(string),
        Version: d.Get("version").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/keycloak/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), strings.NewReader(string(rb)))

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

	var mcloudKeycloakResponse McloudKeycloakResponse
	err = json.Unmarshal(body, &mcloudKeycloakResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("admin_password", mcloudKeycloakResponse.AdminPassword)
    d.Set("fqdn", mcloudKeycloakResponse.Fqdn)
    d.Set("name", mcloudKeycloakResponse.Name)
    d.Set("pki_ca_id", mcloudKeycloakResponse.PkiCaId)
    d.Set("secret_key", mcloudKeycloakResponse.SecretKey)
    d.Set("server_pool_id", mcloudKeycloakResponse.ServerPoolId)
    d.Set("sku", mcloudKeycloakResponse.Sku)
    d.Set("status", mcloudKeycloakResponse.Status)
    d.Set("themes", mcloudKeycloakResponse.Themes)
    d.Set("version", mcloudKeycloakResponse.Version)

	return diags
}

func resourceMcloudKeycloakRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/keycloak/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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
		log.Printf("[WARN] mcloud_keycloak %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var mcloudKeycloakResponse McloudKeycloakResponse
	err = json.Unmarshal(body, &mcloudKeycloakResponse)
	//err = json.NewDecoder(resp.Body).Decode(McloudKeycloakResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("admin_password", mcloudKeycloakResponse.AdminPassword)
    d.Set("fqdn", mcloudKeycloakResponse.Fqdn)
    d.Set("name", mcloudKeycloakResponse.Name)
    d.Set("pki_ca_id", mcloudKeycloakResponse.PkiCaId)
    d.Set("secret_key", mcloudKeycloakResponse.SecretKey)
    d.Set("server_pool_id", mcloudKeycloakResponse.ServerPoolId)
    d.Set("sku", mcloudKeycloakResponse.Sku)
    d.Set("status", mcloudKeycloakResponse.Status)
    d.Set("themes", mcloudKeycloakResponse.Themes)
    d.Set("version", mcloudKeycloakResponse.Version)

	return diags
}
func resourceMcloudKeycloakUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMcloudKeycloakCreate(ctx, d, m)
}

func resourceMcloudKeycloakDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/keycloak/%s", strings.Trim(provider.HostURL, "/"), d.Get("name").(string)), nil)
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