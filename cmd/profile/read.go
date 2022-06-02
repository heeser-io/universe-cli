package profile

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
		Short: "reads the authenticated user profile",
		Run: func(cmd *cobra.Command, args []string) {
			p, err := client.Client.Profile.Read(&v1.ReadProfileParams{})
			if err != nil {
				color.Red(err.Error())
			}
			fmt.Println(string(v1.StructToByte(p)))
		},
	}
)
