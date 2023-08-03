package profile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	AccountID         string
	ClientID          string
	ReadSubprofileCmd = &cobra.Command{
		Use:   "read-subprofile",
		Short: "reads a subprofile",
		Run: func(cmd *cobra.Command, args []string) {
			p, err := client.Client.Profile.ReadSubprofile(&v2.ReadSubprofileParams{
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
	ReadSubprofileCmd.Flags().StringVar(&AccountID, "account-id", "", "accountId of the sub profile")
	ReadSubprofileCmd.Flags().StringVar(&ClientID, "client-id", "", "clientId of the sub profile")
}
