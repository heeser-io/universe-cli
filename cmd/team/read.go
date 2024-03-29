package team

import (
	"fmt"

	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	ReadCmd = &cobra.Command{
		Use:   "read",
		Short: "reads a team with the given team id",
		Run: func(cmd *cobra.Command, args []string) {
			team, err := client.Client.Team.Read(&v1.ReadTeamParams{
				TeamID: TeamID,
			})
			if err != nil {
				panic(err)
			}

			fmt.Println(string(v1.StructToByte(team)))
		},
	}
)

func init() {
}
