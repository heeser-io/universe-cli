package data

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	ListIndicesCmd = &cobra.Command{
		Use:   "list-indices",
		Short: "list all indices for the given collection",
		Run: func(cmd *cobra.Command, args []string) {
			indices, err := client.Client.Data.ListIndices(&v2.ListIndicesParams{
				CollectionName: CollectionName,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(indices)))
		},
	}
)

func init() {
}
