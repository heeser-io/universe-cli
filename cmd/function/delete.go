package function

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "deletes a function with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Function.Delete(&v1.DeleteFunctionParams{
				FunctionID: FunctionID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			color.Green("successfully deleted function %s", FunctionID)
		},
	}
)
