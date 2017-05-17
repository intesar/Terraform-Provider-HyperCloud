package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	//"fmt"
	"log"
	//"github.com/hashicorp/terraform/helper/logging"
	//"strings"
)

func resourceCompute() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeCreate,
		Read:   resourceComputeRead,
		Update: resourceComputeUpdate,
		Delete: resourceComputeDelete,

		Schema: map[string]*schema.Schema{
			"Blueprint-ID" : &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"Name" : &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"IP" : &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"Status" : &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"Create_Date" : &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

		},
	}
}

func resourceComputeCreate(d *schema.ResourceData, m interface{}) error {
	blueprintId := d.Get("Blueprint-ID").(string)
	log.Printf("[HC-INFO] Blueprint ID: %s", blueprintId)
	//log.Printf("[DEBUG] Blueprint ID: %s", blueprintId)

	client := m.(*Auth)

	s := client.create(blueprintId)

	log.Printf("[HC-INFO] Setting Compute ID: %s", s.Results.ID)
	d.SetId(s.Results.ID)

	log.Printf("[HC-INFO] Setting Compute Name: %s", s.Results.Name)
	d.Set("Name", s.Results.Name)

	log.Printf("[HC-INFO] Setting Compute IP: %s", s.Results.HostOrIp)
	d.Set("IP", s.Results.HostOrIp)

	log.Printf("[HC-INFO] Setting Compute Status: %s", s.Results.HostOrIp)
	d.Set("Status", s.Results.DockerServerStatus)

	log.Printf("[HC-INFO] Setting Compute IP: %s", s.Results.HostOrIp)
	d.Set("Create_Date", s.Results.CreateDate)

	return nil
}

func resourceComputeRead(d *schema.ResourceData, m interface{}) error {
	//blueprintId := d.Get("Blueprint-ID").(string)

	client := m.(*Auth)
	log.Printf("[HC-INFO] Reading Compute...")
	log.Printf("[HC-INFO] Compute ID: %s", d.Id())

	s, _ := client.getVM(d.Id())

	log.Printf("[HC-INFO] Compute [%s]...", s)

	status := s.Results.DockerServerStatus
	// Attempt to read from an upstream API
	//obj, ok := client.Get(d.Id())
	log.Printf("[HC-INFO] Compute Status: %s", status)
	ok := (status == "CONNECTED" || status == "PROVISIONED")

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(s.Results.ID)
	d.Set("Name", s.Results.Name)
	d.Set("IP", s.Results.HostOrIp)
	return nil
}

func resourceComputeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[HC-INFO] Updating Compute...")
	log.Printf("[HC-INFO] Compute ID: %s", d.Get("ID"))
	return nil
}

func resourceComputeDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Auth)
	log.Printf("[HC-INFO] Deleting Compute...")
	log.Printf("[HC-INFO] Compute ID: %s", d.Id())
	s, _ := client.delete(d.Id())

	if !s.Errors {
		d.SetId("")
		d.Set("Name", "")
		d.Set("IP", "")
	}
	return nil
}