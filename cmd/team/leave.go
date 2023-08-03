package team

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	LeaveTeamCmd = &cobra.Command{
		Use:   "leave",
		Short: "leaves a team with the given team id",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Team.LeaveTeam(&v2.LeaveTeamParams{
				TeamID: TeamID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			color.Green("successfully left team %s", TeamID)
		},
	}
)

func init() {
}
