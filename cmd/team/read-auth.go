package team

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	ReadAuthCmd = &cobra.Command{
		Use:   "read-auth",
		Short: "reads the current authenticated team",
		Run: func(cmd *cobra.Command, args []string) {
			team, err := client.Client.Team.ReadAuth(&v2.ReadAuthTeamParams{})
			if err != nil {
				color.Red("no team authenticated")
				return
			}

			fmt.Println(string(v2.StructToByte(team)))
		},
	}
)

func init() {
}
