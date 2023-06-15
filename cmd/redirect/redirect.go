package redirect

import (
	"github.com/spf13/cobra"
)

var (
	RedirectID  string
	RedirectCmd = &cobra.Command{
		Use:   "redirect",
		Short: "redirect api",
	}
)

func init() {
	RedirectCmd.PersistentFlags().StringVarP(&RedirectID, "redirect-id", "p", "", "id of the redirect")
	RedirectCmd.AddCommand(CreateCmd)
	RedirectCmd.AddCommand(ReadCmd)
	RedirectCmd.AddCommand(UpdateCmd)
	RedirectCmd.AddCommand(DeleteCmd)
	RedirectCmd.AddCommand(ListCmd)
}
