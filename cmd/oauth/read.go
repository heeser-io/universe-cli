package oauth

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	ReadCmd = &cobra.Command{
		Use:   "read",
		Short: "reads an oauth app with the given id",
		Run: func(cmd *cobra.Command, args []string) {
			oauthObj, err := client.Client.OAuth.Read(&v1.ReadOAuthParams{
				OAuthID: OAuthID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(oauthObj)))
		},
	}
)

func init() {
}
