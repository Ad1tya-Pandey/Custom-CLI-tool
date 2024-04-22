package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactive mode",
	Long:  `Interactive mode allows you to navigate through different commands and options using a menu-based interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a prompt for selecting commands
		prompt := promptui.Select{
			Label: "Select Command",
			Items: []string{"CPU Usage", "Memory Usage", "Disk Usage", "Network Usage", "Process Info", "System Stats"},
		}

		// Prompt the user to select a command
		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		// Execute the selected command
		switch result {
		case "CPU Usage":
			cpuCmd.Run(cmd, args)
		case "Memory Usage":
			memCmd.Run(cmd, args)
		case "Disk Usage":
			diskCmd.Run(cmd, args)
		case "Network Usage":
			networkCmd.Run(cmd, args)
		case "Process Info":
			processCmd.Run(cmd, args)
		case "System Stats":
			statsCmd.Run(cmd, args)
		}
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
	interactiveCmd.Flags().StringVarP(&exportFilePath, "export", "e", "", "Export to file (provide file path)")
}
