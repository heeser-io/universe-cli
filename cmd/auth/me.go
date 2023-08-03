package auth

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	MeCmd = &cobra.Command{
		Use:   "me",
		Short: "shows auth parameters for the current authenticated entity",
		Run: func(cmd *cobra.Command, args []string) {
			authParams, err := client.Client.Auth.Me()
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(authParams)))
		},
	}
)
