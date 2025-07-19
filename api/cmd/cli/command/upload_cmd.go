package command

import "github.com/spf13/cobra"

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload any csv file from identity vendor",
	Long: `Sight is tool for uploading new records from identity vendors, aml screening and lookup for people.
	Pass flags to the program you are debugging using ` + "`--`" + `
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
