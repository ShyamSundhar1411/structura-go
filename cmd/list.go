/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available architectures",
	Long: `The 'list' command displays all available backend architectures 
	that can be used to generate project folder structures. This helps you 
	select the appropriate architecture for your Go project setup.`,
	Run: func(cmd *cobra.Command, args []string) {
		architectures, err := loadAllTemplates("./templates/")
		if err != nil {
			fmt.Println("❌ Error loading templates:", err)
			return
		}
		fmt.Println("📌 Available Architectures:")
		for _, tmpl := range architectures {
			fmt.Printf("🔹 %s: %s\n", tmpl.Architecture, tmpl.Description)
			printFolderStructure(tmpl.Folders,"")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
