/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Project struct {
	Name string
	Path string
	Description string
	Architecture string
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Go project structure",
	Long: `A tool to quickly initialize a new Go project with common architectures
	like MVC, MVCS, etc. It creates the necessary directories based on your
	chosen architecture.`,
	Run: func(cmd *cobra.Command, args []string) {
		var project Project
		project.Name = interactivePrompt("Project Name", "my-project")
		project.Path = interactivePrompt("Project Path", "./"+project.Name)
		project.Description = interactivePrompt("Project Description", "A new Go project")
		project.Architecture = selectArchitecture()
		fmt.Print("Project '%s' initialized successfully at '%s'\n", project.Name, project.Path)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().StringP("name", "n", "", "Name of the project")
	initCmd.Flags().StringP("path", "p", "", "Path to initialize the project")
	initCmd.Flags().StringP("description", "d", "", "Description of the project")
	initCmd.Flags().StringP("architecture", "a", "", "Architecture to use (MVC, MVCS, etc.)")
}
func interactivePrompt(label, defaultValue string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	result, err := prompt.Run()
	if err != nil{
		fmt.Println("Error during generation",err)
		os.Exit(1)
	}
	return result
}

func selectArchitecture()string{
	prompt := promptui.Select{
		Label: "Select Architecture",
		Items: []string{"MVC", "MVCS", "MVP", "MVPF", "MVI"},
	}
	_, result, err := prompt.Run()
	if err != nil{
		fmt.Println("Error during generation", err)
		os.Exit(1)
	}
	return result
}