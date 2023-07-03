package oauth

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "deletes an oauth app with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.OAuth.Delete(&v1.DeleteOAuthParams{
				OAuthID: OAuthID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			color.Green("successfully deleted redirect %s", OAuthID)
		},
	}
)
