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
	"strconv"
	"strings"
)

func resourceServerPoolHcloud() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerPoolHcloudCreate,
		ReadContext:   resourceServerPoolHcloudRead,
		UpdateContext: resourceServerPoolHcloudUpdate,
		DeleteContext: resourceServerPoolHcloudDelete,
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
			"instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceServerPoolHcloudCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	instance_type := d.Get("instance_type").(string)
	instance_count := d.Get("instance_count").(int)

	server_pool := ServerPoolHcloudRequest{
		Name:          name,
		InstanceCount: instance_count,
		InstanceType:  instance_type,
		RunSetup:      true,
	}

	rb, err := json.Marshal(server_pool)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/server-pool/hcloud/%s", strings.Trim(c.HostURL, "/"), name), strings.NewReader(string(rb)))
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

	var serverPoolHcloudResponse ServerPoolHcloudResponse
	err = json.Unmarshal(body, &serverPoolHcloudResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	if !serverPoolHcloudResponse.Success {
		return diag.FromErr(fmt.Errorf(serverPoolHcloudResponse.Error))
	}

	debug(strconv.Itoa(serverPoolHcloudResponse.Task.Id))
	c.waitForTaskToFinish(serverPoolHcloudResponse.Task.Id)

	d.SetId(name)
	d.Set("instance_type", serverPoolHcloudResponse.ServerPoolHcloud.InstanceType)
	d.Set("instance_count", serverPoolHcloudResponse.ServerPoolHcloud.InstanceCount)

	return diags
}

func resourceServerPoolHcloudRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	client := c.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/server-pool/hcloud/%s", strings.Trim(c.HostURL, "/"), name), nil)
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

	var serverPoolHcloudResponse ServerPoolHcloudResponse
	err = json.Unmarshal(body, &serverPoolHcloudResponse)
	//err = json.NewDecoder(resp.Body).Decode(sshKeyResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	if !serverPoolHcloudResponse.Success {
		return diag.FromErr(fmt.Errorf(serverPoolHcloudResponse.Error))
	}

	d.Set("instance_type", serverPoolHcloudResponse.ServerPoolHcloud.InstanceType)
	d.Set("instance_count", serverPoolHcloudResponse.ServerPoolHcloud.InstanceCount)

	return diags
}

func resourceServerPoolHcloudUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceServerPoolHcloudCreate(ctx, d, m)
}

func resourceServerPoolHcloudDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
