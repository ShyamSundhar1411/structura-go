package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ShyamSundhar1411/structura-go/domain"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)



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

func loadAllTemplates(dir string) ([]domain.Template, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var templates []domain.Template
	for _,file := range files{
		if (!strings.HasSuffix(file.Name(), ".yaml") || file.Name() == "initial_structure.yaml"){
			continue
		}
		data,err := os.ReadFile(dir + "/" + file.Name())
		if err != nil {
			fmt.Println("⚠️ Error reading:", file.Name())
			continue
		}
		var tmpl domain.Template
		err = yaml.Unmarshal(data, &tmpl)
		if err != nil {
			fmt.Println("⚠️ Error unmarshalling:", file.Name())
			continue
		}
		templates = append(templates, tmpl)
	}
	return templates,nil
}
func loadTemplateFromArchitecture(dir string, architecture string)(*domain.Template, error){
	filePath := dir+"/"+strings.ToLower(architecture)+".yaml"
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("⚠️ Error reading:", filePath)
		return nil,err
	}
	var template domain.Template
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		fmt.Println("⚠️ Error unmarshalling:", filePath)
		return nil,err
	}
	return &template, nil	
}

func printFolderStructure(structure interface{}, indent string) {
	switch folders := structure.(type) {
	case []interface{}:
		for _, folder := range folders {
			if str, ok := folder.(string); ok {
				fmt.Println(indent + "📂 " + str)
			} else {
				fmt.Println(indent + "⚠️ Unexpected non-string folder:", folder)
			}
		}
	case map[string]interface{}: 
		for parent, subfolders := range folders {
			fmt.Println(indent + "📂 " + parent)
			printFolderStructure(subfolders, indent+"   ") 
		}
	default:
		fmt.Println(indent + "⚠️ Unknown folder structure format:", structure)
	}
}

func assignProjectAttributes(project *domain.Project,cmd *cobra.Command)(*domain.Project){
	orderedFlags := []string{"name", "path", "description", "architecture"}
	attributes := map[string]domain.Attribute{
		"name" : {
			Field: &project.Name,
			Label: "Project Name",
		},
		"path" : {
			Field: &project.Path,
			Label: "Project Path",
		},
		"description" : {
			Field: &project.Description,
			Label: "Project Description",
		},
		"architecture" : {
			Field: &project.Architecture,
			Label: "Project Architecture",
		},
	}
	defaults := map[string] string{
		"name" : "cmd",
		"path" : "./",
		"description" : "A new Go project",
		"architecture" : "MVC",
	}
	architectureOptions :=  []string{"MVC", "MVC-API", "MVCS", "Hexagonal"}
	for _,flag := range orderedFlags{
		attr := attributes[flag]
		if cmd.Flags().Changed(flag){
			value,_ := cmd.Flags().GetString(flag)
			*attr.Field = value
		}else{
			if flag == "architecture"{
				*attr.Field = selectPrompt(attr.Label,architectureOptions)
			}else{
				*attr.Field = interactivePrompt(attr.Label, defaults[flag])
			}
		}
	}
	return project
}
func createFolder(parentPath string,folders interface{})error{
	switch folder := folders.(type) {
	case []interface{}:
		for _, folder := range folder {
			if str, ok := folder.(string); ok {
				folderPath := filepath.Join(parentPath, str)
				err := os.Mkdir(folderPath, 0755)
				if err != nil {
					return fmt.Errorf("⚠️ Error creating folder: %s -> %v", folderPath, err)
					
				}
			} else {
				return fmt.Errorf("⚠️ Unexpected non-string folder: %v", folder)
			}
		}
	case map[string]interface{}:
		for parent, subfolders := range folder {
			folderPath := filepath.Join(parentPath, parent)
			err := os.Mkdir(folderPath, 0755)
			if err != nil {
				return fmt.Errorf("⚠️ Error creating folder: %s -> %v", folderPath, err)
				
			}
			if err := createFolder(folderPath, subfolders); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("⚠️ Unknown folder structure format: %v", folder)
	}
	return nil
}
func createArchitectureStructure(project *domain.Project){
	architecture := strings.ToLower(project.Architecture)
	template, err := loadTemplateFromArchitecture("./templates", architecture)
	if err != nil {
		fmt.Println("⚠️ Error loading template:", err)
		os.Exit(1)
	}
	projectRoot := filepath.Join(project.Path, project.Name)
	if err := os.MkdirAll(projectRoot, 0755); err != nil {
		fmt.Println("⚠️ Error creating project root:", err)
		return
	}
	content := []byte(`package main

import "fmt"

func main() {
	fmt.Println("Hello, Go Project!")
}`)
	filePath := filepath.Join(projectRoot, "main.go")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		fmt.Println("❌ Error writing file:", err)
		return
	}
	fmt.Println("✅ main.go created successfully at", filePath)

	if err := createFolder(project.Path, template.Folders); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("✅ Folder structure created successfully at", projectRoot)
	
}