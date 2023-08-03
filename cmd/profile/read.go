package profile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	ReadCmd = &cobra.Command{
		Use:   "read",
		Short: "reads the authenticated user profile",
		Run: func(cmd *cobra.Command, args []string) {
			p, err := client.Client.Profile.Read(&v2.ReadProfileParams{})
			if err != nil {
				color.Red(err.Error())
			}
			fmt.Println(string(v2.StructToByte(p)))
		},
	}
)
