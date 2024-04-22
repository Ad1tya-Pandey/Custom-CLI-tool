package cmd

import (
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Display network usage statistics",
	Long:  `Display network usage statistics such as bandwidth usage and network interfaces.`,
	Run: func(cmd *cobra.Command, args []string) {
		DisplayNetworkUsage()
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)
	networkCmd.Flags().StringVarP(&exportFilePath, "export", "e", "", "Export to file (provide file path)")

}
