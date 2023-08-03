package data

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "deletes a data object with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Data.Delete(&v2.DeleteDataParams{
				CollectionName: CollectionName,
				IndexName:      IndexName,
				IndexValue:     IndexValue,
				Filter:         Filter,
				AuthField:      AuthField,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			color.Green("successfully deleted data %s.%s", IndexName, IndexValue)
		},
	}
)
