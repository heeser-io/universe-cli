package team

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	Name      string
	Slug      string
	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "creates a team with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			team, err := client.Client.Team.Create(&v2.CreateTeamParams{
				Name: Name,
				Slug: Slug,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			fmt.Println(string(v2.StructToByte(team)))
		},
	}
)

func init() {
}
