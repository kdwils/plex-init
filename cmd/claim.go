package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// claimCmd represents the claim command
var claimCmd = &cobra.Command{
	Use:   "claim",
	Short: "claim something",
	Long:  `claim something`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("claim called")
	},
}

func init() {
	rootCmd.AddCommand(claimCmd)
}
