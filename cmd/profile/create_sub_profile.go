package profile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	Email               string
	Groups              *[]string
	CreateSubprofileCmd = &cobra.Command{
		Use:   "create-subprofile",
		Short: "create a sub profile with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			p, err := client.Client.Profile.CreateSubprofile(&v2.CreateSubprofileParams{
				AccountID: AccountID,
				ClientID:  ClientID,
				Email:     Email,
				Groups:    *Groups,
				// TODO UserMetadata: UserMetadata,
			})
			if err != nil {
				color.Red(err.Error())
			}
			fmt.Println(string(v2.StructToByte(p)))
		},
	}
)

func init() {
	CreateSubprofileCmd.Flags().StringVar(&Email, "email", "", "email of the new subprofile")
	Groups = CreateSubprofileCmd.Flags().StringSlice("groups", nil, "groups of the new subprofile")
	CreateSubprofileCmd.Flags().StringVar(&AccountID, "account-id", "", "accountId of the sub profile")
	CreateSubprofileCmd.Flags().StringVar(&ClientID, "client-id", "", "clientId of the sub profile")
}
