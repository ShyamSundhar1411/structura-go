package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v3"
)

type Template struct {
	Architecture string `yaml:"architecture"`
	Description string `yaml:"description"`
	Folders interface {} `yaml:"folders"`
}

func selectPrompt(label string, options []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: options,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Error during prompt", err)
		os.Exit(1)
	}
	return result
}

func interactivePrompt(label, defaultValue string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Println("Error during prompt", err)
		os.Exit(1)
	}
	return result
}

func loadAllTemplates(dir string) ([]Template, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var templates []Template
	for _,file := range files{
		if (!strings.HasSuffix(file.Name(), ".yaml") || file.Name() == "initial_structure.yaml"){
			continue
		}
		data,err := os.ReadFile(dir + "/" + file.Name())
		if err != nil {
			fmt.Println("⚠️ Error reading:", file.Name())
			continue
		}
		var tmpl Template
		err = yaml.Unmarshal(data, &tmpl)
		if err != nil {
			fmt.Println("⚠️ Error unmarshalling:", file.Name())
			continue
		}
		templates = append(templates, tmpl)
	}
	return templates,nil
}
func loadTemplateFromArchitecture(dir string, architecture string)(*Template, error){
	filePath := dir+"/"+strings.ToLower(architecture)+".yaml"
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("⚠️ Error reading:", filePath)
		return nil,err
	}
	var template Template
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		fmt.Println("⚠️ Error unmarshalling:", filePath)
		return nil,err
	}
	return &template, nil	
}

func printFolderStructure(structure interface{}, indent string) {
	switch folders := structure.(type) {
	case []interface{}: // Flat structure (YAML may return []interface{})
		for _, folder := range folders {
			if str, ok := folder.(string); ok {
				fmt.Println(indent + "📂 " + str)
			} else {
				fmt.Println(indent + "⚠️ Unexpected non-string folder:", folder)
			}
		}
	case map[string]interface{}: // Nested structure
		for parent, subfolders := range folders {
			fmt.Println(indent + "📂 " + parent)
			printFolderStructure(subfolders, indent+"   ") // Recursive call
		}
	default:
		fmt.Println(indent + "⚠️ Unknown folder structure format:", structure)
	}
}

