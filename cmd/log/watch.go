package log

import (
	"github.com/spf13/cobra"
)

var (
	Resources *[]string
	WatchCmd  = &cobra.Command{
		Use:   "watch",
		Short: "watch live logs for the given resources",
		Run: func(cmd *cobra.Command, args []string) {
			// 	logs, err := client.Client.Log.List(&v1.ListLogParams{

			// 	})
			// 	if err != nil {
			// 		color.Red("err:%v\n", err)
			// 	}
			// 	fmt.Println(string(v1.StructToByte(logs)))
		},
	}
)

func init() {
	Resources = WatchCmd.Flags().StringSlice("resources", nil, "resources to watch logs for")
}
