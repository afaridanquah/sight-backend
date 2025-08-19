package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sight",
	Short: "Sight is a cli tool for managing identity",
	Long: `Sight is tool for uploading new records from identity vendors, aml screening and lookup for people.
		Pass flags to the program you are debugging using ` + "`--`" + `
		`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute(config Config) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rootCmd.AddCommand(uploadCmd)
}
