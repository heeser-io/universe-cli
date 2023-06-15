package redirect

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	Target     string
	CreateTags *[]string
	CreateCmd  = &cobra.Command{
		Use:   "create",
		Short: "creates a redirect with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			redirectObj, err := client.Client.Redirect.Create(&v1.CreateRedirectParams{
				Target: Target,
				Tags:   *CreateTags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			fmt.Println(string(v1.StructToByte(redirectObj)))
		},
	}
)

func init() {
	CreateCmd.Flags().StringVar(&Target, "target", "", "url target")
	CreateTags = CreateCmd.Flags().StringSlice("tags", nil, "tags of the redirect")
}
