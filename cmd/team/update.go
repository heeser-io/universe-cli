package team

import (
	"fmt"

	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	UpdateName string
	UpdateSlug string
	UpdateCmd  = &cobra.Command{
		Use:   "update",
		Short: "updates a team with the given team id",
		Run: func(cmd *cobra.Command, args []string) {
			team, err := client.Client.Team.Update(&v1.UpdateTeamParams{
				TeamID: TeamID,
				Name:   UpdateName,
				Slug:   UpdateSlug,
			})
			if err != nil {
				panic(err)
			}

			fmt.Println(string(v1.StructToByte(team)))
		},
	}
)

func init() {
	UpdateCmd.Flags().StringVar(&UpdateName, "name", "", "new name of the team")
	UpdateCmd.Flags().StringVar(&UpdateSlug, "slug", "", "new slug of the team")
}
