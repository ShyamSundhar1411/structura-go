package domain
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

)

func AssignProjectAttributes(project *Project, cmd *cobra.Command) *Project {
	orderedFlags := []string{"name", "path", "description", "architecture","env",".gitignore","readme"}
	attributes := map[string]Attribute{
		"name": {
			Field: &project.Name,
			Label: "Project Name",
		},
		"path": {
			Field: &project.Path,
			Label: "Project Path",
		},
		"description": {
			Field: &project.Description,
			Label: "Project Description",
		},
		"architecture": {
			Field: &project.Architecture,
			Label: "Project Architecture",
		},
		"env" :{
			Field: &project.GenerateEnv,
			Label: "Do you want to generate .env? [y/n]",
		},
		".gitignore":{
			Field: &project.GenerateGitIgnore,
			Label: "Do you want to generate .gitignore? [y/n]",
		},
		"readme":{
			Field: &project.GenerateReadME,
			Label: "Do you want to generate README.md ? [y/n]",
		},
	}
	defaults := map[string]string{
		"name":         "cmd",
		"path":         "./",
		"description":  "A new Go project",
		"architecture": "MVC",
	}
	architectureOptions := []string{"MVC", "MVC-API", "MVCS", "Hexagonal"}
	for _, flag := range orderedFlags {
		attr := attributes[flag]
		if cmd.Flags().Changed(flag) {
			value, _ := cmd.Flags().GetString(flag)
			attr.Field = value
		} else {
			if flag == "architecture" {
				attr.Field = SelectPrompt(attr.Label, architectureOptions)
			} else {

				attr.Field = InteractivePrompt(attr.Label, defaults[flag])
			}
		}
	}
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
	
	if err := CreateBoilerPlates(project); err !=nil{
		fmt.Println(err)
		return
	}
	if err := CreateFolder(project.Path, template.Folders); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("✅ Folder structure created successfully at", projectRoot)

}
func LoadTemplateFromArchitecture(dir string, architecture string)(*Template, error){
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