package commercetools

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/labd/commercetools-go-sdk/commercetools"
	"github.com/labd/commercetools-go-sdk/service/project"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
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
			"shipping_rate_input_type": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"label": {
										Type:     schema.TypeMap,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	log.Fatal("A project can not be created through terraform")
	return fmt.Errorf("A project can not be created through terraform")
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	log.Print("Reading project from commercetools")
	svc := getProjectService(m)

	project, err := svc.Get()

	if err != nil {
		if ctErr, ok := err.(commercetools.Error); ok {
			if ctErr.Code() == commercetools.ErrResourceNotFound {
				return nil
			}
		}
		return err
	}

	log.Print("Found the following project:")
	log.Print(stringFormatObject(project))

	d.SetId(project.Key)
	d.Set("version", project.Version)
	d.Set("name", project.Name)
	d.Set("currencies", project.Currencies)
	d.Set("countries", project.Countries)
	d.Set("languages", project.Languages)
	d.Set("messages", project.Messages)
	d.Set("shippingRateInputType", project.ShippingRateInputType)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	svc := getProjectService(m)

	input := &project.UpdateInput{
		Version: d.Get("version").(int),
		Actions: commercetools.UpdateActions{},
	}

	if d.HasChange("name") {
		input.Actions = append(input.Actions, &project.ChangeName{d.Get("name").(string)})
	}

	if d.HasChange("currencies") {
		newCurrencies := getStringSlice(d, "currencies")
		input.Actions = append(
			input.Actions,
			&project.ChangeCurrencies{Currencies: newCurrencies})
	}

	if d.HasChange("countries") {
		newCountries := getStringSlice(d, "countries")
		input.Actions = append(
			input.Actions,
			&project.ChangeCountries{Countries: newCountries})
	}

	if d.HasChange("languages") {
		newLanguages := getStringSlice(d, "languages")
		input.Actions = append(
			input.Actions,
			&project.ChangeLanguages{Languages: newLanguages})
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
			&project.ChangeMessagesEnabled{MessagesEnabled: enabled})
	}

	if d.HasChange("shipping_rate_input_type") {
		log.Println("=== SHIPPING ===")

		shippingRateInputType, err := resourceProjectGetShippingRateInputType(d)
		if err != nil {
			return err
		}
		log.Println(shippingRateInputType)

		input.Actions = append(
			input.Actions,
			&project.SetShippingRateInputType{ShippingRateInputType: shippingRateInputType})
	}

	_, err := svc.Update(input)
	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	log.Fatal("A project can not be deleted through terraform")
	return fmt.Errorf("A project can not be deleted through terraform")
}

func getProjectService(m interface{}) *project.Service {
	client := m.(*commercetools.Client)
	svc := project.New(client)
	return svc
}

func getStringSlice(d *schema.ResourceData, field string) []string {
	input := d.Get(field).([]interface{})
	var currencyObjects []string
	for _, raw := range input {
		currencyObjects = append(currencyObjects, raw.(string))
	}

	return currencyObjects
}

func resourceProjectGetShippingRateInputType(d *schema.ResourceData) (project.ShippingRateInputType, error) {
	inputType := d.Get("shipping_rate_input_type").(map[string]interface{})
	if inputType != nil {
		switch inputType["type"] {
		case "CartValue":
			return project.CartValue{}, nil
		case "CartScore":
			return project.CartScore{}, nil
		case "CartClassification":
			values := inputType["values"].(map[string]interface{})
			var localizedEnumValues []commercetools.LocalizedEnumValue
			for key, label := range values {
				newEnumValue := commercetools.LocalizedEnumValue{Key: key, Label: label.(map[string]string)}
				localizedEnumValues = append(localizedEnumValues, newEnumValue)
			}
			return project.CartClassification{Values: localizedEnumValues}, nil
		default:
			return nil, fmt.Errorf("ShippingRateInputType %s not implemented", inputType["type"])
		}
	} else {
		return nil, nil
	}
}
