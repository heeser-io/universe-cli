package redirect

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "deletes a redirect with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Redirect.Delete(&v1.DeleteRedirectParams{
				RedirectID: RedirectID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			color.Green("successfully deleted redirect %s", RedirectID)
		},
	}
)
