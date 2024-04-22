package cmd

import (
	"github.com/spf13/cobra"
)

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Display information about running processes",
	Long:  `Display information about running processes such as process ID, CPU and memory usage, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		DisplayProcessInfo()
	},
}

func init() {
	rootCmd.AddCommand(processCmd)
	processCmd.Flags().StringVarP(&exportFilePath, "export", "e", "", "Export to file (provide file path)")

}
