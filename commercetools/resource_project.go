package commercetools

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/labd/commercetools-go-sdk/commercetools"
)

func resourceProjectSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Exists: resourceProjectExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"currencies": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"countries": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"languages": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"messages": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"external_oauth": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"authorization_header": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"shipping_rate_input_type": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"label": {
										Type:     TypeLocalizedString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceProjectExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := getClient(m)

	_, err := client.ProjectGet()
	if err != nil {
		if ctErr, ok := err.(commercetools.ErrorResponse); ok {
			if ctErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := getClient(m)
	project, err := client.ProjectGet()

	if err != nil {
		if ctErr, ok := err.(commercetools.ErrorResponse); ok {
			if ctErr.StatusCode == 404 {
				return nil
			}
		}
		return err
	}

	err = projectUpdate(d, client, project.Version)
	if err != nil {
		return err
	}
	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	log.Print("[DEBUG] Reading projects from commercetools")
	client := getClient(m)

	project, err := client.ProjectGet()

	if err != nil {
		if ctErr, ok := err.(commercetools.ErrorResponse); ok {
			if ctErr.StatusCode == 404 {
				return nil
			}
		}
		return err
	}

	log.Print("[DEBUG] Found the following project:")
	log.Print(stringFormatObject(project))

	d.SetId(project.Key)
	d.Set("version", project.Version)
	d.Set("name", project.Name)
	d.Set("currencies", project.Currencies)
	d.Set("countries", project.Countries)
	d.Set("languages", project.Languages)
	d.Set("shipping_rate_input_type", project.ShippingRateInputType)
	d.Set("external_oauth", project.ExternalOAuth)
	// d.Set("createdAt", project.CreatedAt)
	// d.Set("trialUntil", project.TrialUntil)
	log.Print("[DEBUG] Logging messages enabled")
	log.Print(stringFormatObject(project.Messages))
	d.Set("messages", project.Messages)
	log.Print(stringFormatObject(d))
	// d.Set("shippingRateInputType", project.ShippingRateInputType)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := getClient(m)
	version := d.Get("version").(int)
	err := projectUpdate(d, client, version)
	if err != nil {
		return err
	}
	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func projectUpdate(d *schema.ResourceData, client *commercetools.Client, version int) error {
	input := &commercetools.ProjectUpdateInput{
		Version: version,
		Actions: []commercetools.ProjectUpdateAction{},
	}

	if d.HasChange("name") {
		input.Actions = append(input.Actions, &commercetools.ProjectChangeNameAction{Name: d.Get("name").(string)})
	}

	if d.HasChange("currencies") {
		newCurrencies := []commercetools.CurrencyCode{}
		for _, item := range getStringSlice(d, "currencies") {
			newCurrencies = append(newCurrencies, commercetools.CurrencyCode(item))
		}

		input.Actions = append(
			input.Actions,
			&commercetools.ProjectChangeCurrenciesAction{Currencies: newCurrencies})
	}

	if d.HasChange("countries") {
		newCountries := []commercetools.CountryCode{}
		for _, item := range getStringSlice(d, "countries") {
			newCountries = append(newCountries, commercetools.CountryCode(item))
		}

		input.Actions = append(
			input.Actions,
			&commercetools.ProjectChangeCountriesAction{Countries: newCountries})
	}

	if d.HasChange("languages") {
		newLanguages := []commercetools.Locale{}
		for _, item := range getStringSlice(d, "languages") {
			newLanguages = append(newLanguages, commercetools.Locale(item))
		}
		input.Actions = append(
			input.Actions,
			&commercetools.ProjectChangeLanguagesAction{Languages: newLanguages})
	}

	if d.HasChange("messages") {
		messages := d.Get("messages").(map[string]interface{})
		// ¯\_(ツ)_/¯
		enabled := false
		if messages["enabled"] == "1" {
			enabled = true
		}

		input.Actions = append(
			input.Actions,
			&commercetools.ProjectChangeMessagesEnabledAction{MessagesEnabled: enabled})
	}

	if d.HasChange("shipping_rate_input_type") {
		newShippingRateInputType := d.Get("shipping_rate_input_type").(commercetools.ShippingRateInputType)
		input.Actions = append(
			input.Actions,
			&commercetools.ProjectSetShippingRateInputTypeAction{ShippingRateInputType: newShippingRateInputType})
	}

	if d.HasChange("external_oauth") {
		externalOAuth := d.Get("external_oauth").(map[string]interface{})
		if externalOAuth["url"] != nil && externalOAuth["authorization_header"] != nil {
			newExternalOAuth := commercetools.ExternalOAuth{
				URL:                 externalOAuth["url"].(string),
				AuthorizationHeader: externalOAuth["authorization_header"].(string),
			}
			input.Actions = append(
				input.Actions,
				&commercetools.ProjectSetExternalOAuthAction{ExternalOAuth: &newExternalOAuth})
		} else {
			input.Actions = append(input.Actions, &commercetools.ProjectSetExternalOAuthAction{ExternalOAuth: nil})
		}
	}

	_, err := client.ProjectUpdate(input)
	return err
}

func getStringSlice(d *schema.ResourceData, field string) []string {
	input := d.Get(field).([]interface{})
	var currencyObjects []string
	for _, raw := range input {
		currencyObjects = append(currencyObjects, raw.(string))
	}

	return currencyObjects
}
