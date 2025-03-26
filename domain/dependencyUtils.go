package domain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func replacePlaceholders(content string, defaultImports []string, defaultBootstrapSetup []string) string {

	if strings.Contains(content, "{{CUSTOM_IMPORTS}}") {
		importBlock := "import (\n\t" + strings.Join(defaultImports, "\n\t") + "\n)"
		content = strings.ReplaceAll(content, "{{CUSTOM_IMPORTS}}", importBlock)
	}
	if strings.Contains(content, "{{CUSTOM_BOOTSTRAP_SETUP}}") && len(defaultBootstrapSetup) != 0 {
		importBlock := strings.Join(defaultBootstrapSetup, "\n\t")
		content = strings.ReplaceAll(content, "{{CUSTOM_BOOTSTRAP_SETUP}}", importBlock)
	} else {
		content = strings.ReplaceAll(content, "{{CUSTOM_BOOTSTRAP_SETUP}}", "")
	}
	return content
}
func createDependencyFiles(dependency Dependency, project *Project, generateCommands map[string]string) error {
	if dependency.Content == nil {
		fmt.Println("⚠️ No content for dependency:", dependency.Name)
		return nil
	}
	for folder, fileContent := range dependency.Content {
		var dirPath string
		if folder == "root" {
			dirPath = project.Path
		} else if folder == "app" {
			dirPath = filepath.Join(project.Path, project.Name)
		} else {
			dirPath = filepath.Join(project.Path, folder)
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return fmt.Errorf("❌ Failed to create directory %s: %v", dirPath, err)
			}
			fmt.Println("📂 Created/verified folder:", dirPath)
		}
		for fileName, content := range fileContent.Files {
			filePath := filepath.Join(dirPath, fileName)
			defaultImports := []string{`"fmt"`}
			defaultBootstrapSetup := []string{}
			if generateCommands["env"] == "y" {
				defaultImports = append(defaultImports, fmt.Sprintf(`"%s/bootstrap"`, project.PackageName))
				defaultBootstrapSetup = append(defaultBootstrapSetup, fmt.Sprintf(`app:=bootstrap.App()`))
				defaultBootstrapSetup = append(defaultBootstrapSetup, fmt.Sprintf(`env:=app.Env`))
				defaultBootstrapSetup = append(defaultBootstrapSetup, fmt.Sprintf(`fmt.Println(env.AppEnv)`))
			}
			if generateCommands["serverType"] == "gin" {
				defaultImports = append(defaultImports, `"github.com/gin-gonic/gin"`)
			} else if generateCommands["serverType"] == "fiber" {
				defaultImports = append(defaultImports, `"github.com/gofiber/fiber/v2"`)
			} else if generateCommands["serverType"] == "echo" {
				defaultImports = append(defaultImports, `"github.com/labstack/echo/v4"`)
				defaultImports = append(defaultImports, `"net/http"`)
			} else if generateCommands["serverType"] == "chi" {
				defaultImports = append(defaultImports, `"net/http"`)
				defaultImports = append(defaultImports, `"github.com/go-chi/chi/v5"`)
			}
			modifiedContent := replacePlaceholders(content, defaultImports, defaultBootstrapSetup)
			if err := os.WriteFile(filePath, []byte(modifiedContent), 0644); err != nil {
				return fmt.Errorf("❌ Failed to create file %s: %v", filePath, err)
			}
			fmt.Println("📄 Created/verified file:", filePath)
		}
	}
	return nil
}

func InstallDependencyPackages(project *Project, generateCommands map[string]string) error {
	projectRoot := filepath.Join(project.Path, project.Name)
	dependencies := project.Dependencies
	for _, dependency := range dependencies {
		cmd := exec.Command("go", "get", dependency.Source)
		cmd.Dir = projectRoot
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("⚠️ Failed to install dependency %s: %v", dependency.Source, err)
		}
		if dependency.Content != nil {
			createDependencyFiles(dependency, project, generateCommands)
		}
	}
	return nil
}
