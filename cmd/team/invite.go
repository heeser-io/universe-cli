package team

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	InviteEmail  string
	InviteAccess string
	InviteCmd    = &cobra.Command{
		Use:   "invite",
		Short: "invites a user to the team",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Team.InviteMember(&v2.InviteMemberParams{
				TeamID: TeamID,
				Email:  InviteEmail,
				Access: InviteAccess,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			color.Green("successfully invited user %s", InviteEmail)
		},
	}
)

func init() {
	InviteCmd.Flags().StringVar(&InviteEmail, "email", "", "e-mail address of the user you want to invite")
	InviteCmd.Flags().StringVar(&InviteAccess, "access", "read", "read or write access")
}
