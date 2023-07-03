package oauth

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	UpdateTarget string
	UpdateTags   *[]string
	UpdateCmd    = &cobra.Command{
		Use:   "update",
		Short: "updates an oauth app with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			oauthObj, err := client.Client.OAuth.Update(&v1.UpdateOAuthParams{
				OAuthID: OAuthID,
				Tags:    *UpdateTags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(oauthObj)))
		},
	}
)

func init() {
	UpdateTags = UpdateCmd.Flags().StringSlice("tags", nil, "new tags of the oauth app")
}
