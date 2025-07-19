package command

import (
	"github.com/spf13/cobra"
)

func New(cfg Config) *cobra.Command {
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

	// repo := postgres.New(pool, pool)

	// service, _ := customerbus.New(repo, logger)

	return cmd
}
