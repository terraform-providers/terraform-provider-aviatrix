package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixBranchRouterAvxTgwAttachment_basic(t *testing.T) {
	if os.Getenv("SKIP_BRANCH_ROUTER_AVX_TGW_ATTACHMENT") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_BRANCH_ROUTER_AVX_TGW_ATTACHMENT is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_branch_router_avx_tgw_attachment.test_branch_router_avx_tgw_attachment"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAviatrixBranchRouterAvxTgwAttachmentPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBranchRouterAvxTgwAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchRouterAvxTgwAttachmentNoOptions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBranchRouterAvxTgwAttachmentExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"backup_pre_shared_key",
					"pre_shared_key",
					"backup_local_tunnel_ip",
					"backup_remote_tunnel_ip",
					"local_tunnel_ip",
					"remote_tunnel_ip",
				},
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAviatrixBranchRouterAvxTgwAttachmentPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBranchRouterAvxTgwAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchRouterAvxTgwAttachmentBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBranchRouterAvxTgwAttachmentExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"backup_pre_shared_key",
					"pre_shared_key",
					"backup_local_tunnel_ip",
					"backup_remote_tunnel_ip",
				},
			},
		},
	})
}

func testAccBranchRouterAvxTgwAttachmentBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_branch_router_avx_tgw_attachment" "test_branch_router_avx_tgw_attachment" {
	branch_name               = "%s"
	transit_gateway_name      = "%s"
	connection_name           = "connection-%s"
	transit_gateway_bgp_asn   = 65000
	branch_router_bgp_asn     = 65001
	phase1_authentication     = "SHA-256"
	phase1_dh_groups          = 14
	phase1_encryption         = "AES-256-CBC"
	phase2_authentication     = "HMAC-SHA-256"
	phase2_dh_groups          = 14
	phase2_encryption         = "AES-256-CBC"
	enable_global_accelerator = true
	enable_branch_router_ha   = false
	pre_shared_key            = "key"
	local_tunnel_ip           = "10.0.0.1/30"
	remote_tunnel_ip          = "10.0.0.2/30"
}
`, os.Getenv("BRANCH_ROUTER_NAME"), os.Getenv("TRANSIT_GATEWAY_NAME"), rName)
}

func testAccBranchRouterAvxTgwAttachmentNoOptions(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_branch_router_avx_tgw_attachment" "test_branch_router_avx_tgw_attachment" {
	branch_name               = "%s"
	transit_gateway_name      = "%s"
	connection_name           = "connection-noopts-%s"
	transit_gateway_bgp_asn   = 65000
	branch_router_bgp_asn     = 65001

}
`, os.Getenv("BRANCH_ROUTER_NAME"), os.Getenv("TRANSIT_GATEWAY_NAME"), rName)
}

func testAccCheckBranchRouterAvxTgwAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("branch_router_avx_tgw_attachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no branch_router_avx_tgw_attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		brata := &goaviatrix.BranchRouterAvxTgwAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		_, err := client.GetBranchRouterAvxTgwAttachment(brata)
		if err != nil {
			return err
		}
		if brata.ConnectionName != rs.Primary.ID {
			return fmt.Errorf("branch_router_avx_tgw_attachment not found")
		}

		return nil
	}
}

func testAccCheckBranchRouterAvxTgwAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_branch_router_avx_tgw_attachment" {
			continue
		}
		foundBranchRouterAvxTgwAttachment := &goaviatrix.BranchRouterAvxTgwAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}
		_, err := client.GetBranchRouterAvxTgwAttachment(foundBranchRouterAvxTgwAttachment)
		if err == nil {
			return fmt.Errorf("branch_router_avx_tgw_attachment still exists")
		}
	}

	return nil
}

func testAccAviatrixBranchRouterAvxTgwAttachmentPreCheck(t *testing.T) {
	if os.Getenv("BRANCH_ROUTER_NAME") == "" {
		t.Fatal("BRANCH_ROUTER_NAME must be set for aviatrix_branch_router_avx_tgw_attachment acceptance test.")
	}
	if os.Getenv("TRANSIT_GATEWAY_NAME") == "" {
		t.Fatal("TRANSIT_GATEWAY_NAME must be set for aviatrix_branch_router_avx_tgw_attachment acceptance test.")
	}
}
