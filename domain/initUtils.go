package domain

import (
	"embed"
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
func LoadDependencies(fs embed.FS,filePath string) ([]Dependency, error) {
	data, err := fs.ReadFile(filePath)
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
func LoadDependency(fs embed.FS,filePath string) (Dependency, error) {
	data, err := fs.ReadFile(filePath)
	if err != nil {
		fmt.Println("⚠️ Error reading:", filePath)
		return Dependency{}, err
	}
	var dependency Dependency
	err = yaml.Unmarshal(data, &dependency)
	if err != nil {
		fmt.Println("⚠️ Error unmarshalling:", filePath)
		return Dependency{}, err
	}
	return dependency, nil
}
func AssignProjectAttributes(project *Project, cmd *cobra.Command) (*Project, map[string]string) {
	orderedFlags := []string{"name", "package-name", "path", "description", "architecture", "env", "generate-server", "server"}
	var generateServer, server, generateEnv string
	var dependencies []Dependency

	attributes := map[string]Attribute{
		"name": {
			Field:        &project.Name,
			Label:        "App Name",
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

		if flag == "env" && generateEnv == "y" {
			filePath := fmt.Sprintf("templates/%s.yaml","default_dependencies")
			defaultDependencies, err := LoadDependencies(TemplateFS,filePath)
			if err != nil {
				fmt.Println("⚠️ Error loading default dependencies:", err)
			}
			for _, dependency := range defaultDependencies {
				dependencies = append(dependencies, dependency)
			}
		}
		if flag == "server" && generateServer == "y" {
			filePath := fmt.Sprintf("templates/%s_server.yaml", server)
			serverDependency, err := LoadDependency(TemplateFS,filePath)
			if err != nil {
				fmt.Println("⚠️ Error loading server dependencies:", err)
			}
			dependencies = append(dependencies, serverDependency)
		}
	}
	generateCommands := map[string]string{
		"env":        generateEnv,
		"server":     generateServer,
		"serverType": server,
	}
	project.PackageName = GetDefaultPackageName() + "/" + project.Name
	project.Dependencies = dependencies
	return project, generateCommands
}

func CreateArchitectureStructure(project *Project, generateComands map[string]string) {
	architecture := strings.ToLower(project.Architecture)
	
	template, err := LoadTemplateFromArchitecture(TemplateFS, architecture)
	if err != nil {
		fmt.Println("⚠️ Error loading template:", err)
		return
	}
	appRoot := filepath.Join(project.Path, project.Name)
	if err := os.MkdirAll(appRoot, 0755); err != nil {
		fmt.Println("⚠️ Error creating project root:", err)
		return
	}

	if err := CreateBoilerPlates(TemplateFS,project); err != nil {
		fmt.Println(err)
		return
	}
	if err := runInitCommands(project.Path, project.PackageName); err != nil {
		fmt.Println("⚠️ Error initializing Go module:", err)
		return
	}
	if err := InstallDependencyPackages(project, generateComands); err != nil {
		fmt.Println("⚠️ Error installing dependency packages:", err)
		return
	}
	if err := CreateFolder(project.Path, template.Folders); err != nil {
		fmt.Println(err)
		return
	}
	if err := CreateFile("readme.MD",template.Readme,project.Path); err != nil{
		fmt.Println(err)
		return
	}
	if err := formatCommand(project.Path);err != nil{
		fmt.Println("⚠️ Error formatting code:", err)
		return
	}
	fmt.Println("✅ Folder structure created successfully at app root", appRoot)

}
func LoadTemplateFromArchitecture(fs embed.FS, architecture string) (*Template, error) {
	filePath := fmt.Sprintf("templates/%s.yaml", strings.ToLower(architecture))
	data, err := fs.ReadFile(filePath)
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

func runInitCommands(projectPath, moduleName string) error {
	fmt.Printf("🔹 Initializing Go module in %s...\n", moduleName)
	goModInitCmd := exec.Command("go", "mod", "init", moduleName)
	goModInitCmd.Dir = projectPath
	goModInitCmd.Stdout = os.Stdout
	goModInitCmd.Stderr = os.Stderr
	if err := goModInitCmd.Run(); err != nil {
		return fmt.Errorf("⚠️ Failed to initialize Go module: %w", err)
	}
	fmt.Println("✅ Go module initialized successfully!")

	fmt.Printf("🔹 Initializing Git repository in %s...\n", projectPath)
	gitInitCmd := exec.Command("git", "init")
	gitInitCmd.Dir = projectPath
	gitInitCmd.Stdout = os.Stdout
	gitInitCmd.Stderr = os.Stderr
	if err := gitInitCmd.Run(); err != nil {
		return fmt.Errorf("⚠️ Failed to initialize Git repository: %w", err)
	}
	fmt.Println("✅ Git repository initialized successfully!")

	return nil
}

func formatCommand(projectPath string) error{
	goFmtCmd := exec.Command("gofmt","-w","-s",".")
	goFmtCmd.Dir = projectPath
	goFmtCmd.Stdout = os.Stdout
	goFmtCmd.Stderr = os.Stderr
	return goFmtCmd.Run()
}
