package domain

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func LoadAllTemplates(dir string) ([]Template, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var templates []Template
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yaml") || file.Name() == "initial_structure.yaml" {
			continue
		}
		data, err := os.ReadFile(dir + "/" + file.Name())
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
	return templates, nil
}

func PrintFolderStructure(structure interface{}, indent string) {
	switch folders := structure.(type) {
	case []interface{}:
		for _, folder := range folders {
			if str, ok := folder.(string); ok {
				fmt.Println(indent + "📂 " + str)
			} else {
				fmt.Println(indent+"⚠️ Unexpected non-string folder:", folder)
			}
		}
	case map[string]interface{}:
		for parent, subfolders := range folders {
			fmt.Println(indent + "📂 " + parent)
			PrintFolderStructure(subfolders, indent+"   ")
		}
	default:
		fmt.Println(indent+"⚠️ Unknown folder structure format:", structure)
	}
}
