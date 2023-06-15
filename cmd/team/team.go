package team

import (
	"github.com/spf13/cobra"
)

var (
	TeamID  string
	TeamCmd = &cobra.Command{
		Use:   "team",
		Short: "team api",
	}
)

func init() {
	TeamCmd.AddCommand(ReadAuthCmd)
	TeamCmd.AddCommand(CreateCmd)
	TeamCmd.AddCommand(InviteCmd)
	TeamCmd.AddCommand(ReadCmd)
	TeamCmd.AddCommand(UpdateCmd)
	TeamCmd.AddCommand(DeleteCmd)
	TeamCmd.AddCommand(ListCmd)
	TeamCmd.AddCommand(ConfirmInviteCmd)
	TeamCmd.AddCommand(LeaveTeamCmd)
	TeamCmd.AddCommand(UpdateMemberCmd)
	TeamCmd.AddCommand(RemoveMemberCmd)

	TeamCmd.PersistentFlags().StringVar(&TeamID, "team-id", "", "id of the team")
}
