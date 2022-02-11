package mcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMcloudOrderBasic(t *testing.T) {
	coffeeID := "1"
	quantity := "2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMcloudOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMcloudOrderConfigBasic(coffeeID, quantity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMcloudOrderExists("mcloud_order.new"),
				),
			},
		},
	})
}

func testAccCheckMcloudOrderDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mcloud_order" {
			continue
		}

		orderID := rs.Primary.ID

		err := c.DeleteOrder(orderID)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckMcloudOrderConfigBasic(coffeeID, quantity string) string {
	return fmt.Sprintf(`
	resource "mcloud_order" "new" {
		items {
			coffee {
				id = %s
			}
    		quantity = %s
  		}
	}
	`, coffeeID, quantity)
}

func testAccCheckMcloudOrderExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OrderID set")
		}

		return nil
	}
}
