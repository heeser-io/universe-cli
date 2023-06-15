package team

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	UpdateMemberAccess string
	UpdateMemberEmail  string
	UpdateMemberCmd    = &cobra.Command{
		Use:   "update-member",
		Short: "updates a team member with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Team.UpdateMember(&v1.UpdateMemberParams{
				TeamID: TeamID,
				Email:  UpdateMemberEmail,
				Access: UpdateMemberAccess,
			})
			if err != nil {
				panic(err)
			}
			color.Green("successfully updated member %s with access: %s", UpdateMemberEmail, UpdateMemberAccess)
		},
	}
)

func init() {
	UpdateMemberCmd.Flags().StringVar(&UpdateMemberEmail, "email", "", "email of the member you want to update")
	UpdateMemberCmd.Flags().StringVar(&UpdateMemberAccess, "access", "", "read or write")
}
