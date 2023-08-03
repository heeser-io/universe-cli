package redirect

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	UpdateTarget string
	UpdateTags   *[]string
	UpdateCmd    = &cobra.Command{
		Use:   "update",
		Short: "updates a redirect with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			redirectObj, err := client.Client.Redirect.Update(&v2.UpdateRedirectParams{
				RedirectID: RedirectID,
				Tags:       *UpdateTags,
				Target:     UpdateTarget,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(redirectObj)))
		},
	}
)

func init() {
	UpdateCmd.Flags().StringVar(&UpdateTarget, "target", "", "new target")
	UpdateTags = UpdateCmd.Flags().StringSlice("tags", nil, "new tags of the redirect")
}
