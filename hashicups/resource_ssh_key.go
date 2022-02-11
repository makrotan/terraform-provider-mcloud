package hashicups

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

func resourceSshKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSshKeyCreate,
		ReadContext:   resourceSshKeyRead,
		UpdateContext: resourceSshKeyUpdate,
		DeleteContext: resourceSshKeyDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: true,
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceSshKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	public_key := d.Get("public_key").(string)
	private_key := d.Get("private_key").(string)

	ssh_key := SshKey{
		Name:       name,
		PublicKey:  public_key,
		PrivateKey: private_key,
	}

	rb, err := json.Marshal(ssh_key)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(c.HostURL, "/"), name), strings.NewReader(string(rb)))
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

	var sshKeyResponse SshKeyResponse
	err = json.Unmarshal(body, &sshKeyResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	if !sshKeyResponse.Success {
		return diag.FromErr(fmt.Errorf(sshKeyResponse.Error))
	}

	d.SetId(name)
	d.Set("public_key", sshKeyResponse.SshKey.PublicKey)
	d.Set("private_key", sshKeyResponse.SshKey.PrivateKey)

	return diags
}

func resourceSshKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	client := c.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(c.HostURL, "/"), name), nil)
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
		log.Printf("[WARN] ssh key %s not present", name)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var sshKeyResponse SshKeyResponse
	err = json.Unmarshal(body, &sshKeyResponse)
	//err = json.NewDecoder(resp.Body).Decode(sshKeyResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	if !sshKeyResponse.Success {
		return diag.FromErr(fmt.Errorf(sshKeyResponse.Error))
	}

	d.Set("public_key", sshKeyResponse.SshKey.PublicKey)
	d.Set("private_key", sshKeyResponse.SshKey.PrivateKey)

	return diags
}

func resourceSshKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSshKeyCreate(ctx, d, m)
}

func resourceSshKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	client := c.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(c.HostURL, "/"), name), nil)
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

//func flattenOrderItems(orderItems *[]OrderItem) []interface{} {
//	if orderItems != nil {
//		ois := make([]interface{}, len(*orderItems), len(*orderItems))
//
//		for i, orderItem := range *orderItems {
//			oi := make(map[string]interface{})
//
//			oi["coffee"] = flattenCoffee(orderItem.Coffee)
//			oi["quantity"] = orderItem.Quantity
//			ois[i] = oi
//		}
//
//		return ois
//	}
//
//	return make([]interface{}, 0)
//}
//
//func flattenCoffee(coffee Coffee) []interface{} {
//	c := make(map[string]interface{})
//	c["id"] = coffee.ID
//	c["name"] = coffee.Name
//	c["teaser"] = coffee.Teaser
//	c["description"] = coffee.Description
//	c["price"] = coffee.Price
//	c["image"] = coffee.Image
//
//	return []interface{}{c}
//}
