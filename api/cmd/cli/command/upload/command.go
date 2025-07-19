package upload

import (
	"bitbucket.org/msafaridanquah/verifylab-service/api/cmd/cli/command"
	"github.com/spf13/cobra"
)

func New(cfg command.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sight",
		Short: "Sight is a cli tool for managing identity",
		Long: `Sight is tool for uploading new records from identity vendors, aml screening and lookup for people.
		Pass flags to the program you are debugging using ` + "`--`" + `
		`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	return cmd
}
