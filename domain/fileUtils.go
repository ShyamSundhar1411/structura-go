package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

func CreateFolder(parentPath string, folders interface{}) error {
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
			if err := CreateFolder(folderPath, subfolders); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("⚠️ Unknown folder structure format: %v", folder)
	}
	return nil
}
func CreateBoilerPlates(projectRoot string) error{
boilerPlateFilePath := filepath.Join("templates", "initial_structure.yaml")
	data, err := os.ReadFile(boilerPlateFilePath)
	if err != nil {
		return fmt.Errorf("❌ Error loading boiler plates: %v", err)
	}
	var boilerPlates []BoilerPlate
	err = yaml.Unmarshal(data, &boilerPlates)
	if err != nil {
		return fmt.Errorf("❌ Error unmarshalling boiler plates: %v", err)
	}
	for _, file := range boilerPlates {
		content := []byte(file.Content)
		filePath := filepath.Join(projectRoot, file.Name)
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return fmt.Errorf("❌ Error writing file %s:%v", file.Name,err)
		}
		fmt.Println("✅ ", file.Name, "created successfully at", filePath)

	}
	return nil
}


