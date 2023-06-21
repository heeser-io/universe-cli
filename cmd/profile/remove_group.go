package profile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	RemoveGroupCmd = &cobra.Command{
		Use:   "remove-group",
		Short: "remove a group from root or sub account",
		Run: func(cmd *cobra.Command, args []string) {
			p, err := client.Client.Profile.RemoveGroup(&v1.RemoveGroupParams{
				Group:     Group,
				AccountID: AccountID,
				ClientID:  ClientID,
			})
			if err != nil {
				color.Red(err.Error())
			}
			fmt.Println(string(v1.StructToByte(p)))
		},
	}
)

func init() {
	RemoveGroupCmd.Flags().StringVar(&Group, "group", "", "group to add")
	RemoveGroupCmd.Flags().StringVar(&AccountID, "account-id", "", "accountId of the sub profile")
	RemoveGroupCmd.Flags().StringVar(&ClientID, "client-id", "", "clientId of the sub profile")
}
