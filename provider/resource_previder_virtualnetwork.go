package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/previder/previder-go-sdk/client"
	"log"
	"time"
)

func resourcePreviderVirtualNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourcePreviderVirtualNetworkCreate,
		Read:   resourcePreviderVirtualNetworkRead,
		Delete: resourcePreviderVirtualNetworkDelete,
		Update: resourcePreviderVirtualNetworkUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcePreviderVirtualNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.BaseClient)

	// Build up our creation options
	network := &client.VirtualNetworkUpdate{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	log.Printf("[DEBUG] VirtualNetwork create configuration: %#v", network)
	log.Printf("[DEBUG] VirtualNetwork create configuration: %#v", c)
	task, err := c.VirtualNetwork.Create(network)

	if err != nil {
		return fmt.Errorf("Error creating VirtualNetwork: %s", err)
	}

	c.Task.WaitFor(task.Id, 5*time.Minute)
	log.Printf("[INFO] Virtual network %s created", network.Name)

	return resourcePreviderVirtualNetworkUpdate(d, meta)
}

func resourcePreviderVirtualNetworkRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.BaseClient)

	log.Printf("[INFO] Retrieving network with name: %s", d.Get("name").(string))
	virtualNetwork, err := c.VirtualNetwork.Get(d.Get("name").(string))
	if err != nil {
		if err.(*client.ApiError).Code == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving VirtualNetwork: %s", err)
	}
	d.SetId(virtualNetwork.Id)
	d.Set("name", virtualNetwork.Name)
	//d.Set("addressPool", virtualNetwork.AddressPool)

	return nil
}

func resourcePreviderVirtualNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.BaseClient)

	log.Printf("[INFO] Retrieving network with name: %s", d.Get("name").(string))
	virtualNetwork, err := c.VirtualNetwork.Get(d.Get("name").(string))
	if err != nil {
		if err.(*client.ApiError).Code == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving VirtualNetwork: %s", err)
	}

	// Build up addressPool
	var update client.VirtualNetworkUpdate
	//
	//if d.Get("address_pool").(map[string]interface{})["ip_start"] != nil {
	//	var addressPool client.AddressPool
	//	addressPool.Start = d.Get("address_pool").(map[string]interface{})["ip_start"].(string)
	//	addressPool.End = d.Get("address_pool").(map[string]interface{})["ip_end"].(string)
	//	addressPool.Mask = d.Get("address_pool").(map[string]interface{})["ip_netmask"].(string)
	//	addressPool.Gateway = d.Get("address_pool").(map[string]interface{})["ip_gateway"].(string)
	//	addressPool.NameServers = append(addressPool.NameServers, d.Get("address_pool").(map[string]interface{})["ip_nameserver1"].(string), d.Get("address_pool").(map[string]interface{})["ip_nameserver2"].(string))
	//	update.AddressPool = addressPool
	//}

	//update.Id = virtualNetwork.Id
	update.Name = virtualNetwork.Name
	//update.PublicNet = virtualNetwork.PublicNet
	//update.VlanId = virtualNetwork.VlanId
	//update.Type = virtualNetwork.Type

	_, uerr := c.VirtualNetwork.Update(virtualNetwork.Id, &update)
	//
	if uerr != nil {
		return fmt.Errorf("Error updating VirtualNetwork: %s", uerr)
	}
	//
	log.Printf("[INFO] Virtual network %s updated", update.Name)

	return resourcePreviderVirtualNetworkRead(d, meta)
}

func resourcePreviderVirtualNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.BaseClient)

	log.Printf("[INFO] Deleting VirtualNetwork: %s", d.Id())
	task, err := c.VirtualNetwork.Delete(d.Id())
	c.Task.WaitFor(task.Id, 5*time.Minute)

	if err != nil {
		return fmt.Errorf("Error deleting VirtualNetwork: %s", err)
	}
	d.SetId("")

	return nil
}
