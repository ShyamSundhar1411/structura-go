/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/ShyamSundhar1411/structura-go/domain"
	"github.com/spf13/cobra"
)





// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Go project structure",
	Long: `A tool to quickly initialize a new Go project with common architectures
	like MVC, MVCS, etc. It creates the necessary directories based on your
	chosen architecture.`,
	Run: func(cmd *cobra.Command, args []string) {
		var project domain.Project
		
		project = *domain.AssignProjectAttributes(&project,cmd)
		
		template,err := domain.LoadTemplateFromArchitecture("./templates/",project.Architecture)
		if err != nil {
			fmt.Println("Error loading template:", err)
			return
		}
		
		fmt.Printf("\n📂 Project Structure for '%s' (%s Architecture):\n", project.Name, project.Architecture)
		domain.PrintFolderStructure(template.Folders, "")
		domain.CreateArchitectureStructure(&project)
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
	initCmd.Flags().StringP("env","e","","Include Env file")
}
