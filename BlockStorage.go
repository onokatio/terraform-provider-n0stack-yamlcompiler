package main

import (
	"fmt"
	"os"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"github.com/n0stack/n0stack/n0proto.go/pkg/dag"
	"gopkg.in/yaml.v3"
)

func resource_n0stack_blockstorage() *schema.Resource {
	return &schema.Resource{
		Create: resource_n0stack_blockstorage_create,
		Read: resource_n0stack_blockstorage_read,
		Update: resource_n0stack_blockstorage_update,
		Delete: resource_n0stack_blockstorage_delete,

		Schema: map[string]*schema.Schema{
			"image_name": {
				Type:	 schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:	 schema.TypeString,
				Optional: true,
			},
			"blockstorage_name": {
				Type:	 schema.TypeString,
				Required: true,
			},
			"annotations": {
				Type:	 schema.TypeMap,
				Required: true,
				Elem: schema.TypeString,
			},
			"labels": {
				Type:	 schema.TypeMap,
				Optional: true,
				Elem: schema.TypeString,
			},
			"request_bytes": {
				Type:	 schema.TypeInt,
				Required: true,
			},
			"limit_bytes": {
				Type:	 schema.TypeInt,
				Required: true,
			},
			"source_url": {
				Type:	 schema.TypeString,
				Optional: true,
			},
		},
	}
}


type BlockStorageTask struct {
	Type        string      `yaml:"type"`
	Action      string      `yaml:"action"`
	Args        struct {
		ImageName           string             `yaml:"image_name"`
		Tag                 string             `yaml:"tag"`
		BlockStorageName    string             `yaml:"block_storage_names"`
		Annotations         map[string]string  `yaml:"annotations"`
		Labels              map[string]string  `yaml:"labels"`
		RequestBytes        uint64             `yaml:"request_bytes"`
		LimitBytes          uint64             `yaml:"limit_bytes"`
		SourceUrl           string             `yaml:"source_url"`
	}
	DependsOn   []string    `yaml:"depends_on"`
	IgnoreError bool        `yaml:"ignore_error"`
	// Rollback []*Task `yaml:"rollback"`

	child   []string
	depends int
}

func resource_n0stack_blockstorage_create(d *schema.ResourceData, meta interface{}) error {

	task := BlockStorageTask{}
	task.Type = "Image"
	task.Action = "GenerateBlockStorage"
	task.Args.ImageName = d.Get("image_name").(string)
	task.Args.Tag = d.Get("tag").(string)
	task.Args.BlockStorageName = d.Get("blockstorage_name").(string)
	task.Args.Annotations = interfaceMap2stringMap(d.Get("annotations").(map[string]interface{}))
	task.Args.Labels = interfaceMap2stringMap(d.Get("labels").(map[string]interface{}))
	task.Args.RequestBytes = uint64(d.Get("request_bytes").(int))
	task.Args.LimitBytes = uint64(d.Get("limit_bytes").(int))
	task.Args.SourceUrl = d.Get("source_url").(string)

	taskList := make(map[string]BlockStorageTask)
	taskList["GenerateBlockStorage-" + d.Get("blockstorage_name").(string)] = task

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

	file, err := os.Create("n0cli-yaml/Generate/BlockStorage-" + d.Get("blockstorage_name").(string) + ".yaml")
	if err != nil {
		return err;
	}

	fmt.Fprint(file, string(yamlString))

	d.SetId(d.Get("blockstorage_name").(string))

	return resource_n0stack_blockstorage_read(d, meta)
}

func resource_n0stack_blockstorage_read(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resource_n0stack_blockstorage_update(d *schema.ResourceData, meta interface{}) error {
	return resource_n0stack_blockstorage_read(d, meta)
}

func resource_n0stack_blockstorage_delete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
