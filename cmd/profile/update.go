package profile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	Avatar   string
	Metadata map[string]string

	UpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "update a root or sub account",
		Run: func(cmd *cobra.Command, args []string) {
			p, err := client.Client.Profile.Update(&v1.UpdateProfileParams{
				AccountID: AccountID,
				ClientID:  ClientID,
				Avatar:    Avatar,
				Metadata:  Metadata,
			})
			if err != nil {
				color.Red(err.Error())
			}
			fmt.Println(string(v1.StructToByte(p)))
		},
	}
)

func init() {
	UpdateCmd.Flags().StringVar(&AccountID, "account-id", "", "accountId of the sub profile")
	UpdateCmd.Flags().StringVar(&ClientID, "client-id", "", "clientId of the sub profile")
	UpdateCmd.Flags().StringToStringVar(&Metadata, "metadata", nil, "new metadata")
	UpdateCmd.Flags().StringVar(&Avatar, "avatar", "", "new avatar url")
}
