package team

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "deletes a team with the given team id",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Team.Delete(&v1.DeleteTeamParams{
				TeamID: TeamID,
			})
			if err != nil {
				panic(err)
			}

			color.Green("successfully deleted team %s", TeamID)
		},
	}
)
