package oauth

import (
	"github.com/spf13/cobra"
)

var (
	OAuthID  string
	OAuthCmd = &cobra.Command{
		Use:   "oauth",
		Short: "oauth api",
	}
)

func init() {
	OAuthCmd.PersistentFlags().StringVarP(&OAuthID, "oauth-id", "p", "", "id of the oauth app")
	OAuthCmd.AddCommand(CreateCmd)
	OAuthCmd.AddCommand(ReadCmd)
	OAuthCmd.AddCommand(UpdateCmd)
	OAuthCmd.AddCommand(DeleteCmd)
	OAuthCmd.AddCommand(ListCmd)
}
