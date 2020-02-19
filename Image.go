package main

import (
	"fmt"
	"os"
        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/yaml.v3"
)

func resource_n0stack_image() *schema.Resource {
	return &schema.Resource{
		Create: resource_n0stack_image_create,
		Read: resource_n0stack_image_read,
		Update: resource_n0stack_image_update,
		Delete: resource_n0stack_image_delete,

		Schema: map[string]*schema.Schema{
			"image_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"blockstorage_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type ImageTask struct {
	Type        string      `yaml:"type"`
	Action      string      `yaml:"action"`
	Args        struct {
		Name           string             `yaml:"name"`
	}
	DependsOn   []string    `yaml:"depends_on"`
	IgnoreError bool        `yaml:"ignore_error"`
	// Rollback []*Task `yaml:"rollback"`

	child   []string
	depends int
}

type ImageRegistryTask struct {
	Type        string      `yaml:"type"`
	Action      string      `yaml:"action"`
	Args        struct {
		ImageName           string             `yaml:"image_name"`
		Tags                []string           `yaml:"tags"`
		BlockStorageName    string             `yaml:"blockstorage_name"`
	}
	DependsOn   []string    `yaml:"depends_on"`
	IgnoreError bool        `yaml:"ignore_error"`
	// Rollback []*Task `yaml:"rollback"`

	child   []string
	depends int
}

func resource_n0stack_image_create(d *schema.ResourceData, meta interface{}) error {
	task := ImageTask{}
	task.Type = "Image"
	task.Action = "ApplyImage"
	task.Args.Name = d.Get("image_name").(string)

	taskList := make(map[string]ImageTask)
	taskList["ApplyImage-" + d.Get("image_name").(string)] = task

	yamlString, err := yaml.Marshal(&taskList)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	task.Action = "DeleteImage"
	taskList_delete := make(map[string]ImageTask)
	taskList_delete["DeleteImage-" + d.Get("image_name").(string)] = task

	yamlString_delete, err := yaml.Marshal(&taskList_delete)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	task_registry := ImageRegistryTask{}
	task_registry.Type = "Image"
	task_registry.Action = "RegisterBlockStorage"
	task_registry.Args.ImageName = d.Get("image_name").(string)
	task_registry.Args.Tags = interfaceList2stringList(d.Get("tags").([]interface{}))
	task_registry.Args.BlockStorageName = d.Get("blockstorage_name").(string)

	taskList_registry := make(map[string]ImageRegistryTask)
	taskList_registry["ApplyImage-" + d.Get("image_name").(string)] = task_registry

	yamlString_registry, err := yaml.Marshal(&taskList_registry)
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

	file, err := os.Create("n0cli-yaml/Generate/ApplyImage-" + d.Get("image_name").(string) + ".yaml")
	if err != nil {
		return err;
	}

	file_delete, err := os.Create("n0cli-yaml/Delete/DeleteImage-" + d.Get("image_name").(string) + ".yaml")
	if err != nil {
		return err;
	}

	file_registry, err := os.Create("n0cli-yaml/Generate/RegisterBlockStorage-" + d.Get("image_name").(string) + "to" + d.Get("blockstorage_name").(string) + ".yaml")
	if err != nil {
		return err;
	}

	fmt.Fprint(file, string(yamlString))
	fmt.Fprint(file_delete, string(yamlString_delete))
	fmt.Fprint(file_registry, string(yamlString_registry))

	d.SetId(d.Get("image_name").(string))

	return resource_n0stack_image_read(d, meta)
}

func resource_n0stack_image_read(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resource_n0stack_image_update(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resource_n0stack_image_delete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
