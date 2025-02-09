package domain

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func SelectPrompt(label string, options []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: options,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Error during prompt", err)
		os.Exit(1)
	}
	return result
}

func InteractivePrompt(label, defaultValue string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Println("Error during prompt", err)
		os.Exit(1)
	}
	return result
}

func validateInput(input string, validValues []string) error {
	input = strings.ToLower(strings.TrimSpace(input))
	for _, v := range validValues {
		if input == strings.ToLower(v) {
			return nil
		}
	}
	return fmt.Errorf("invalid input, please enter one of: %s", strings.Join(validValues, ", "))
}
