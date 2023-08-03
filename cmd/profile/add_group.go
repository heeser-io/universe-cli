package profile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	Group       string
	AddGroupCmd = &cobra.Command{
		Use:   "add-group",
		Short: "add a group to root or sub account",
		Run: func(cmd *cobra.Command, args []string) {
			p, err := client.Client.Profile.AddGroup(&v2.AddGroupParams{
				Group:     Group,
				AccountID: AccountID,
				ClientID:  ClientID,
			})
			if err != nil {
				color.Red(err.Error())
			}
			fmt.Println(string(v2.StructToByte(p)))
		},
	}
)

func init() {
	AddGroupCmd.Flags().StringVar(&Group, "group", "", "group to add")
	AddGroupCmd.Flags().StringVar(&AccountID, "account-id", "", "accountId of the sub profile")
	AddGroupCmd.Flags().StringVar(&ClientID, "client-id", "", "clientId of the sub profile")
}
