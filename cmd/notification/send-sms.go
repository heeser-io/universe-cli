package notification

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	Phone    string
	Message  string
	SenderID string

	SendSMSCmd = &cobra.Command{
		Use:   "send-sms",
		Short: "sends an sms with the given parameters",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.SMS.SendSMS(&v1.SendSMSParams{
				Phone:    Phone,
				Message:  Message,
				SenderID: SenderID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			color.Green("successfully sent sms to %s", Phone)
		},
	}
)

func init() {
	SendSMSCmd.Flags().StringVar(&Phone, "phone", "", "phone number with country code")
	SendSMSCmd.Flags().StringVar(&Message, "message", "", "message to send")
	SendSMSCmd.Flags().StringVar(&SenderID, "sender-id", "", "custom sender id with max length 11 (defaults to UNIVERSE)")
}
