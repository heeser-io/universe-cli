package notification

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	Receiver   *[]string
	CC         *[]string
	Subject    string
	Body       string
	TemplateID string
	ClientID   string
	Variables  map[string]string

	SendEmailCmd = &cobra.Command{
		Use:   "send-email",
		Short: "sends an email with the given parameters",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Email.SendEmail(&v1.SendEmailParams{
				Receiver:   *Receiver,
				CC:         *CC,
				Subject:    Subject,
				Message:    Body,
				ClientID:   ClientID,
				TemplateID: TemplateID,
				Variables:  Variables,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}

			for _, r := range *Receiver {
				color.Green("successfully sent email to %s", r)
			}
			for _, c := range *CC {
				color.Green("successfully sent email to %s", c)
			}
		},
	}
)

func init() {
	Receiver = SendEmailCmd.Flags().StringSlice("receiver", nil, "multiple receiver separated by comma")
	CC = SendEmailCmd.Flags().StringSlice("cc", nil, "multiple cc separated by comma")
	SendEmailCmd.Flags().StringVar(&Subject, "subject", "", "subject of the email")
	SendEmailCmd.Flags().StringVar(&Body, "body", "", "body of the email")
	SendEmailCmd.Flags().StringVar(&ClientID, "client-id", "", "id of the client to use")
	SendEmailCmd.Flags().StringVar(&TemplateID, "template-id", "", "id of the template to use")
	SendEmailCmd.Flags().StringToStringVar(&Variables, "variables", map[string]string{}, "variables filled inside of a template (e.g name=foo)")
}
