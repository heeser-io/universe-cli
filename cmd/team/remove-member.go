package team

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	RemoveMemberEmail string
	RemoveMemberCmd   = &cobra.Command{
		Use:   "remove-member",
		Short: "removes a team member with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Team.RemoveMember(&v1.RemoveMemberParams{
				TeamID: TeamID,
				Email:  RemoveMemberEmail,
			})
			if err != nil {
				panic(err)
			}
			color.Green("successfully removed member %s", RemoveMemberEmail)
		},
	}
)

func init() {
	RemoveMemberCmd.Flags().StringVar(&RemoveMemberEmail, "email", "", "email of the member you want to remove")
}
