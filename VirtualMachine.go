package main

import (
	"fmt"
	"os"
        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/yaml.v3"
)

func resource_n0stack_virtualmachine() *schema.Resource {
	return &schema.Resource{
		Create: resource_n0stack_virtualmachine_create,
		Read: resource_n0stack_virtualmachine_read,
		Update: resource_n0stack_virtualmachine_update,
		Delete: resource_n0stack_virtualmachine_delete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
			"request_cpu_milli_core": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"limit_cpu_milli_core": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"request_memory_bytes": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"limit_memory_bytes": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"block_storage_names": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ssh_authorized_keys": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"nics": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ipv4_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ipv6_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

type Nics struct {
	NetworkName  string `yaml:"network_name"`
	Ipv4Address string `yaml:"ipv4_address"`
	Ipv6Address string `yaml:"ipv6_address"`
}

type VirtualMachineTask struct {
	Type        string      `yaml:"type"`
	Action      string      `yaml:"action"`
	Args        struct {
		Name           string             `yaml:"name"`
		Tag                 string             `yaml:"tag"`
		Annotations         map[string]string  `yaml:"annotations"`
		Labels              map[string]string  `yaml:"labels"`
		BlockStorageNames   []string             `yaml:"block_storage_names"`
		RequestCpuMilliCore uint32             `yaml:"request_cpu_milli_core"`
		LimitCpuMilliCore   uint32               `yaml:"limit_cpu_milli_core"`
		RequestMemoryBytes  uint64              `yaml:"request_memory_bytes"`
		LimitMemoryBytes    uint64                `yaml:"limit_memory_bytes"`
		SshAuthorizedKeys   []string             `yaml:"ssh_authorized_keys"`
		Nics                []Nics
	}
	DependsOn   []string    `yaml:"depends_on"`
	IgnoreError bool        `yaml:"ignore_error"`
	// Rollback []*Task `yaml:"rollback"`

	child   []string
	depends int
}

func resource_n0stack_virtualmachine_create(d *schema.ResourceData, meta interface{}) error {

	task := VirtualMachineTask{}
	task.Type = "VirtualMachine"
	task.Action = "CreateVirtualMachine"
	task.Args.Name = d.Get("name").(string)
	task.Args.Annotations = interfaceMap2stringMap(d.Get("annotations").(map[string]interface{}))
	task.Args.Labels = interfaceMap2stringMap(d.Get("labels").(map[string]interface{}))
	task.Args.BlockStorageNames = interfaceList2stringList(d.Get("block_storage_names").([]interface{}))
	task.Args.RequestCpuMilliCore = uint32(d.Get("request_cpu_milli_core").(int))
	task.Args.LimitCpuMilliCore = uint32(d.Get("limit_cpu_milli_core").(int))
	task.Args.RequestMemoryBytes = uint64(d.Get("request_memory_bytes").(int))
	task.Args.LimitMemoryBytes = uint64(d.Get("limit_memory_bytes").(int))
	task.Args.SshAuthorizedKeys = interfaceList2stringList(d.Get("ssh_authorized_keys").([]interface{}))
	nics := make([]Nics,0)
	for _, value := range (d.Get("nics").([]interface{})) {
		nic := Nics{}
		element := value.(map[string]interface{})
		nic.NetworkName = element["network_name"].(string)
		nic.Ipv4Address = element["ipv4_address"].(string)
		nic.Ipv6Address = element["ipv6_address"].(string)
		nics = append(nics, nic)
	}
	task.Args.Nics = nics

	taskList := make(map[string]VirtualMachineTask)
	taskList["GenerateVirtualMachine-" + d.Get("name").(string)] = task

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

	file, err := os.Create("n0cli-yaml/Generate/VirtualMachine-" + d.Get("name").(string) + ".yaml")
	if err != nil {
		return err;
	}

	fmt.Fprint(file, string(yamlString))

	d.SetId(d.Get("name").(string))

	return resource_n0stack_virtualmachine_read(d, meta)
}

func resource_n0stack_virtualmachine_read(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resource_n0stack_virtualmachine_update(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resource_n0stack_virtualmachine_delete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
