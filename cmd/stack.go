package cmd

import (
	"github.com/heeser-io/universe-cli/cmd/stack"
	"github.com/spf13/cobra"
)

var (
	StackCmd = &cobra.Command{
		Use: "stack",
	}
)

func init() {
	StackCmd.AddCommand(stack.PushCmd)
	StackCmd.AddCommand(stack.VerifyCmd)
	StackCmd.AddCommand(stack.RemoveCmd)
}
