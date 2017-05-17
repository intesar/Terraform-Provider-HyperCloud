package main

import (
"github.com/hashicorp/terraform/helper/schema"
	"fmt"
)

func Provider() *schema.Provider {
	return &schema.Provider {
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"accessKey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"secretKey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hypercloud_compute": resourceCompute(),
		},
		ConfigureFunc: configure,


	}
}

func configure(d *schema.ResourceData) (interface{}, error) {

	url := d.Get("url").(string)
	accessKey := d.Get("accessKey").(string)
	secretKey := d.Get("secretKey").(string)

	fmt.Printf("Creating HC Client with URL [%s] Key [%s] Secret [%s]", url, accessKey, secretKey)

	client := newAuth(url, accessKey, secretKey)

	return client, nil
}