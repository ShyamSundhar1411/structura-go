package domain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func AssignProjectAttributes(project *Project, cmd *cobra.Command) *Project {
	orderedFlags := []string{"name", "path", "description", "architecture", "env", "server"}
	var server string
	attributes := map[string]Attribute{
		"name": {
			Field: &project.Name,
			Label: "Project Name",
			Type: "Prompt",
		},
		"path": {
			Field: &project.Path,
			Label: "Project Path",
			Type: "Prompt",
		},
		"description": {
			Field: &project.Description,
			Label: "Project Description",
			Type: "Prompt",
		},
		"architecture": {
			Field: &project.Architecture,
			Label: "Project Architecture",
			Type: "Select",
			Options: []string{"MVC", "MVC-API", "MVCS", "Hexagonal"},
		},
		"env": {
			Field: &project.GenerateEnv,
			Label: "Do you want to generate .env? [y/n]",
			Type: "Prompt",
		},
		
		"server":{
			Field: &server,
			Label: "Project Server",
			Type: "Select",
			Options: []string{"gin", "fiber", "echo", "chi", "none"},
		},
	}
	defaults := map[string]string{
		"name":         "cmd",
		"path":         "./",
		"description":  "A new Go project",
		"architecture": "MVC",
	}
	dependencies := []string{}
	for _, flag := range orderedFlags {
		attr := attributes[flag]
		if cmd.Flags().Changed(flag) {
			value, _ := cmd.Flags().GetString(flag)
			*attr.Field = value
		} else {
			if attr.Type == "Select" {
				options, ok := attr.Options.([]string)
				if !ok {
					fmt.Println("Error: Options is not of type []string")
					os.Exit(1)
				}
				*attr.Field = SelectPrompt(attr.Label, options)
				fmt.Println(server)
				if flag == "server"{
					dependencies = append(dependencies, server)
				}
			} else {
				*attr.Field = InteractivePrompt(attr.Label, defaults[flag])
				if flag == "env" {
					dependencies = append(dependencies, "viper")
					
				}
			}
		}
	}
	project.Dependencies = dependencies
	return project
}

func CreateArchitectureStructure(project *Project) {
	architecture := strings.ToLower(project.Architecture)
	template, err := LoadTemplateFromArchitecture("./templates", architecture)
	if err != nil {
		fmt.Println("⚠️ Error loading template:", err)
		os.Exit(1)
	}
	projectRoot := filepath.Join(project.Path, project.Name)
	if err := os.MkdirAll(projectRoot, 0755); err != nil {
		fmt.Println("⚠️ Error creating project root:", err)
		return
	}

	if err := CreateBoilerPlates(project); err != nil {
		fmt.Println(err)
		return
	}
	if err := runGoModInit(project.Path, project.Name); err != nil {
		fmt.Println("⚠️ Error initializing Go module:", err)
		return
	}
	if err := installDependencyPackages(project.Path,project.Dependencies); err != nil {
		fmt.Println("⚠️ Error installing dependency packages:", err)
		return
	}
	if err := CreateFolder(project.Path, template.Folders); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("✅ Folder structure created successfully at", projectRoot)

}
func LoadTemplateFromArchitecture(dir string, architecture string) (*Template, error) {
	filePath := dir + "/" + strings.ToLower(architecture) + ".yaml"
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("⚠️ Error reading:", filePath)
		return nil, err
	}
	var template Template
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		fmt.Println("⚠️ Error unmarshalling:", filePath)
		return nil, err
	}
	return &template, nil
}

func runGoModInit(projectRoot, moduleName string) error {
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func LoadDependencies(filePath string) ([]Dependency, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("⚠️ Error reading:", filePath)
		return nil, err
	}
	var dependencies []Dependency
	err = yaml.Unmarshal(data, &dependencies)
	if err != nil {
		fmt.Println("⚠️ Error unmarshalling:", filePath)
		return nil, err
	}
	return dependencies, nil
}
func GetDependencySource(dependencies []Dependency, name string) (string, error) {
	for _, dep := range dependencies {
		if dep.Name == name {
			return dep.Source, nil
		}
	}
	return "", fmt.Errorf("❌ Dependency '%s' not found", name)
}
func installDependencyPackages(projectRoot string, stringDependencies []string) error {
	filePath := "./templates/default_dependencies.yaml"
	dependencies, err := LoadDependencies(filePath)
	if err != nil {
		return fmt.Errorf("⚠️ Error loading dependencies: %v", err)
	}

	for _, dep := range stringDependencies {
		source, err := GetDependencySource(dependencies, dep)
		if err != nil {
			return fmt.Errorf("⚠️ Error resolving dependency %s: %v", dep, err)
		}
		fmt.Println("Installing:", source)

		cmd := exec.Command("go", "get", source)
		cmd.Dir = projectRoot
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("⚠️ Failed to install dependency %s: %v", source, err)
		}
	}

	return nil
}