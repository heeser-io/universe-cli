package team

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	ConfirmInviteCmd = &cobra.Command{
		Use:   "confirm-invite",
		Short: "confirms an team invitation for the current authenticated user",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Team.ConfirmInvite(&v1.ConfirmInviteParams{
				TeamID: TeamID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			color.Green("successfully confirmed invite for team %s", TeamID)
		},
	}
)

func init() {
}
