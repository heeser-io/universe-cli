package cmd

import (
	"fmt"
	"os"

	"github.com/heeser-io/universe-cli/cmd/auth"
	"github.com/heeser-io/universe-cli/cmd/data"
	"github.com/heeser-io/universe-cli/cmd/file"
	"github.com/heeser-io/universe-cli/cmd/function"
	"github.com/heeser-io/universe-cli/cmd/log"
	"github.com/heeser-io/universe-cli/cmd/notification"
	"github.com/heeser-io/universe-cli/cmd/project"
	"github.com/heeser-io/universe-cli/cmd/quota"
	"github.com/heeser-io/universe-cli/cmd/redirect"
	"github.com/heeser-io/universe-cli/cmd/task"
	"github.com/heeser-io/universe-cli/cmd/team"
	"github.com/heeser-io/universe-cli/cmd/websocket"
	"github.com/spf13/cobra"
)

var (
	TeamID  string
	rootCmd = &cobra.Command{
		Use:           "universe",
		Short:         "CLI for universe",
		SilenceErrors: true,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(StackCmd)
	rootCmd.AddCommand(ProfileCmd)
	rootCmd.AddCommand(team.TeamCmd)
	rootCmd.AddCommand(file.FileCmd)
	rootCmd.AddCommand(function.FunctionCmd)
	rootCmd.AddCommand(log.LogCmd)
	rootCmd.AddCommand(project.ProjectCmd)
	rootCmd.AddCommand(data.DataCmd)
	rootCmd.AddCommand(notification.NotificationCmd)
	rootCmd.AddCommand(redirect.RedirectCmd)
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(quota.QuotaCmd)
	rootCmd.AddCommand(websocket.WebsocketCmd)
	rootCmd.AddCommand(task.TaskCmd)

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "whether to show debug logs or not")
	rootCmd.PersistentFlags().String("branch", "v", "set a branch")
}
