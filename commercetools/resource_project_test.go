package commercetools

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccProjectCreate_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "name", "Test this thing",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "countries.#", "3",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "currencies.#", "2",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "languages.#", "4",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "messages.enabled", "true",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "external_oauth.url", "https://example.com/oauth/token",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "external_oauth.authorization_header", "Bearer secret",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "shipping_rate_input_type.type", "CartValue",
					),
				),
			},
			{
				Config: testAccProjectConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "name", "Test this thing new",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "countries.#", "4",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "currencies.#", "3",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "languages.#", "5",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "messages.enabled", "false",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "external_oauth.url", "https://new-example.com/oauth/token",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "external_oauth.authorization_header", "Bearer new-secret",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "shipping_rate_input_type.type", "CartClassification",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "shipping_rate_input_type.values.#", "3",
					),
				),
			},
			{
				Config: testAccProjectConfigDeleteOAuth(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "name", "Test this thing new",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "countries.#", "4",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "currencies.#", "3",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "languages.#", "5",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "messages.enabled", "false",
					),
					resource.TestCheckNoResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "external_oauth.url",
					),
					resource.TestCheckNoResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "external_oauth.authorization_header",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "shipping_rate_input_type.type", "CartClassification",
					),
					resource.TestCheckResourceAttr(
						"commercetools_project_settings.acctest_project_settings", "shipping_rate_input_type.values.#", "3",
					),
				),
			},
		},
	})
}

func testAccCheckProjectDestroy(s *terraform.State) error {
	return nil
}

func testAccProjectConfig() string {
	return `
		resource "commercetools_project_settings" "acctest_project_settings" {
			name       = "Test this thing"
			countries  = ["NL", "DE", "US"]
			currencies = ["EUR", "USD"]
			languages  = ["nl", "de", "en", "en-US"]
			external_oauth = {
				url = "https://example.com/oauth/token"
				authorization_header = "Bearer secret"
			}
			messages = {
			  enabled = true
			}
			shipping_rate_input_type = { 
				type = "CartValue" 
            }
		}`
}

func testAccProjectConfigUpdate() string {
	return `
		resource "commercetools_project_settings" "acctest_project_settings" {
			name       = "Test this thing new"
			countries  = ["NL", "DE", "US", "GB"]
			currencies = ["EUR", "USD", "GBP"]
			languages  = ["nl", "de", "en", "en-US", "fr"]
			external_oauth = {
				url = "https://new-example.com/oauth/token"
				authorization_header = "Bearer new-secret"
			}
			messages = {
			  enabled = false
			}
            shipping_rate_input_type = { 
				type = "CartClassification"
				values = [ 
					{ 
						key = "Small"
						label = { 
							"en" = "Small",
							"de" = "Klein"
						}
					},
					{ 
						key = "Medium",
						label = { 
							"en" = "Medium",
							"de" = "Mittel" 
						} 	
					},
					{ 
						key = "Heavy",
						label = { 
							"en" = "Heavy",
							"de" = "Schwergut"
						} 
					},
				]
            }
		}`
}

func testAccProjectConfigDeleteOAuth() string {
	return `
		resource "commercetools_project_settings" "acctest_project_settings" {
			name       = "Test this thing new"
			countries  = ["NL", "DE", "US", "GB"]
			currencies = ["EUR", "USD", "GBP"]
			languages  = ["nl", "de", "en", "en-US", "fr"]
			messages = {
			  enabled = false
			}
            shipping_rate_input_type = { 
				type = "CartClassification"
				values = [ 
					{ 
						key = "Small"
						label = { 
							"en" = "Small",
							"de" = "Klein"
						}
					},
					{ 
						key = "Medium",
						label = { 
							"en" = "Medium",
							"de" = "Mittel" 
						} 
					},
					{ 
						key = "Heavy",
						label = { 
							"en" = "Heavy",
							"de" = "Schwergut"
						} 
					},
				]
            }
		}`
}
