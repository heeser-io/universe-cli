package data

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	CountCmd = &cobra.Command{
		Use:   "count",
		Short: "returns the number of elements in a data collection",
		Run: func(cmd *cobra.Command, args []string) {
			countObj, err := client.Client.Data.Count(&v1.CountDataParams{
				CollectionName: CollectionName,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(countObj)))
		},
	}
)

func init() {

}
