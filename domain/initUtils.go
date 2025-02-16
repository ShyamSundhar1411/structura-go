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
func GetGitHubUsername() string {
	cmd := exec.Command("git", "config", "--global", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func GetDefaultPackageName() string {
	username := GetGitHubUsername()
	if username == "" {
		return ""
	}
	return fmt.Sprintf("github.com/%s", username)
}
func AssignProjectAttributes(project *Project, cmd *cobra.Command) *Project {
    orderedFlags := []string{"name", "package-name", "path", "description", "architecture", "env", "generate-server", "server"}
    var generateServer, server,generateEnv string
    dependencies := []string{}

    attributes := map[string]Attribute{
        "name": {
            Field:        &project.Name,
            Label:        "Project Name",
            Type:         "Prompt",
            DefaultValue: "cmd",
        },
        "package-name": {
            Field:        &project.PackageName,
            Label:        "Project Package Name",
            Type:         "Prompt",
            DefaultValue: GetDefaultPackageName(),
        },
        "path": {
            Field:        &project.Path,
            Label:        "Project Path",
            Type:         "Prompt",
            DefaultValue: "./",
        },
        "description": {
            Field:        &project.Description,
            Label:        "Project Description",
            Type:         "Prompt",
            DefaultValue: "A new Go project",
        },
        "architecture": {
            Field:        &project.Architecture,
            Label:        "Project Architecture",
            Type:         "Select",
            Options:      []string{"MVC", "MVC-API", "MVCS", "Hexagonal"},
            DefaultValue: "MVC",
        },
        "env": {
            Field:        &generateEnv,
            Label:        "Do you want to generate .env? [y/n]",
            Type:         "Prompt",
            DefaultValue: "n",
        },
        "generate-server": {
            Field:        &generateServer,
            Label:        "Do you want to generate a server? [y/n]",
            Type:         "Prompt",
            DefaultValue: "n",
        },
        "server": {
            Field:        &server,
            Label:        "Choose the server framework",
            Type:         "Select",
            Options:      []string{"gin", "fiber", "echo", "chi", "none"},
            DefaultValue: "none",
            Condition: func() bool {
                return generateServer == "y"
            },
        },
    }

    for _, flag := range orderedFlags {
        attr := attributes[flag]

        if cmd.Flags().Changed(flag) {
            value, _ := cmd.Flags().GetString(flag)
            *attr.Field = value
        } else {
            if attr.Type == "Select" && (attr.Condition == nil || attr.Condition()) {
                *attr.Field = SelectPrompt(attr.Label, attr.Options.([]string))
            } else if attr.Type == "Prompt" {
                *attr.Field = InteractivePrompt(attr.Label, attr.DefaultValue)
            }
        }

        if flag == "env" && generateEnv == "y"{
            dependencies = append(dependencies, "viper","logrus")
        }
        if flag == "server" && generateServer == "y" {
            dependencies = append(dependencies, server)
        }
    }

    project.PackageName = GetDefaultPackageName() + "/" + project.Name
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
	if err := runGoModInit(project.Path, project.PackageName); err != nil {
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