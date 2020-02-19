package main

import (
	"fmt"
	"os"
	//"google.golang.org/grpc"
        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//ppool "github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"gopkg.in/yaml.v3"
)

func resource_n0stack_network() *schema.Resource {
	return &schema.Resource{
		Create: resource_n0stack_network_create,
		Read: resource_n0stack_network_read,
		Update: resource_n0stack_network_update,
		Delete: resource_n0stack_network_delete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ipv4_cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ipv6_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: schema.TypeString,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: schema.TypeString,
			},
		},
	}
}

type NetworkTask struct {
	Type        string      `yaml:"type"`
	Action      string      `yaml:"action"`
	Args        struct {
		Name           string             `yaml:"name"`
		Annotations         map[string]string  `yaml:"annotations"`
		Labels              map[string]string  `yaml:"labels"`
		Ipv4Cidr           string             `yaml:"ipv4_cidr"`
		Ipv6Cidr           string             `yaml:"ipv6_cidr"`
	}
	DependsOn   []string    `yaml:"depends_on"`
	IgnoreError bool        `yaml:"ignore_error"`
	// Rollback []*Task `yaml:"rollback"`

	child   []string
	depends int
}

func resource_n0stack_network_create(d *schema.ResourceData, meta interface{}) error {
	task := NetworkTask{}
	task.Type = "Network"
	task.Action = "ApplyNetwork"
	task.Args.Name = d.Get("name").(string)
	task.Args.Annotations = interfaceMap2stringMap(d.Get("annotations").(map[string]interface{}))
	task.Args.Labels = interfaceMap2stringMap(d.Get("labels").(map[string]interface{}))
	task.Args.Ipv4Cidr = d.Get("ipv4_cidr").(string)
	task.Args.Ipv6Cidr = d.Get("ipv6_cidr").(string)

	taskList := make(map[string]NetworkTask)
	taskList["ApplyNetwork-" + d.Get("name").(string)] = task

	yamlString, err := yaml.Marshal(&taskList)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	err = os.MkdirAll("n0cli-yaml/Generate", 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll("n0cli-yaml/Delete", 0755)
	if err != nil {
		return err
	}

	file, err := os.Create("n0cli-yaml/Generate/Network-" + d.Get("name").(string) + ".yaml")
	if err != nil {
		return err;
	}

	fmt.Fprint(file, string(yamlString))
	d.SetId(d.Get("name").(string))

	return resource_n0stack_network_read(d, meta)
}

func resource_n0stack_network_read(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resource_n0stack_network_update(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resource_n0stack_network_delete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
