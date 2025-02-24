package domain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)


func createDependencyFiles(dependency Dependency, projectPath string) error {
	if dependency.Content == nil {
		fmt.Println("⚠️ No content for dependency:", dependency.Name)
		return nil
	}
	projectRoot := projectPath
	for folder, fileContent := range dependency.Content {
		var dirPath string
		if folder == "root"{
			dirPath = projectRoot
		} else{
			dirPath =  filepath.Join(projectRoot, folder)
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return fmt.Errorf("❌ Failed to create directory %s: %v", dirPath, err)
			}
			fmt.Println("📂 Created/verified folder:", dirPath)
		}
		for fileName, content := range fileContent.Files {
			filePath := filepath.Join(dirPath, fileName)
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("❌ Failed to create file %s: %v", filePath, err)
			}
			fmt.Println("📄 Created/verified file:", filePath)
		}
	}
	return nil
}

func InstallDependencyPackages(project *Project) error {
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
		if dependency.Content != nil{
			createDependencyFiles(dependency, project.Path)
		}
	}
	return nil
}