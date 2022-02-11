package hashicups

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

	//o, err := c.CreateOrder(ois)
	//if err != nil {
	//	return diag.FromErr(err)
	//}

	//d.SetId(strconv.Itoa(o.ID))
	d.SetId(name)

	resourceSshKeyRead(ctx, d, m)

	return diags
}

func resourceSshKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	orderID := d.Id()

	order, err := c.GetOrder(orderID)
	if err != nil {
		return diag.FromErr(err)
	}

	orderItems := flattenOrderItems(&order.Items)
	if err := d.Set("items", orderItems); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSshKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	orderID := d.Id()

	if d.HasChange("items") {
		items := d.Get("items").([]interface{})
		ois := []OrderItem{}

		for _, item := range items {
			i := item.(map[string]interface{})

			co := i["coffee"].([]interface{})[0]
			coffee := co.(map[string]interface{})

			oi := OrderItem{
				Coffee: Coffee{
					ID: coffee["id"].(int),
				},
				Quantity: i["quantity"].(int),
			}
			ois = append(ois, oi)
		}

		_, err := c.UpdateOrder(orderID, ois)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceOrderRead(ctx, d, m)
}

func resourceSshKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	orderID := d.Id()

	err := c.DeleteOrder(orderID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

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
