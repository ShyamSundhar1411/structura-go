/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ShyamSundhar1411/structura-go/cmd"
)

func main() {
	

	binaryPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}


	newBinaryPath := binaryPath
	if len(binaryPath) > 3 && binaryPath[len(binaryPath)-3:] == "-go" {
		newBinaryPath = binaryPath[:len(binaryPath)-3]
	}


	if binaryPath != newBinaryPath {
		err := exec.Command("mv", binaryPath, newBinaryPath).Run()
		if err != nil {
			fmt.Println("Error renaming binary:", err)
		} else {
			fmt.Println("✅ Installed successfully as 'structura'")
		}
	}
	cmd.Execute()
}
