package data

import (
	"github.com/spf13/cobra"
)

var (
	CollectionName string
	IndexName      string
	IndexValue     string
	AuthField      string

	DataCmd = &cobra.Command{
		Use:   "data",
		Short: "data api",
	}
)

func init() {
	DataCmd.Flags().StringToStringVar(&Filter, "filter", map[string]string{}, "filters")
	DataCmd.PersistentFlags().StringVarP(&CollectionName, "collection-name", "c", "", "name of the collection (required)")
	DataCmd.PersistentFlags().StringVar(&IndexName, "index-name", "_id", "name of the index")
	DataCmd.PersistentFlags().StringVar(&IndexValue, "index-value", "", "value of the index")
	DataCmd.PersistentFlags().StringVar(&AuthField, "auth-field", "", "e.g profileId (default is userId)")

	DataCmd.MarkFlagRequired("collection-name")
	DataCmd.AddCommand(CreateCmd)
	DataCmd.AddCommand(FindOrCreateCmd)
	DataCmd.AddCommand(ReadCmd)
	DataCmd.AddCommand(CountCmd)
	DataCmd.AddCommand(UpdateCmd)
	DataCmd.AddCommand(DeleteCmd)
	DataCmd.AddCommand(ListCmd)
	DataCmd.AddCommand(PushArrayCmd)
	DataCmd.AddCommand(UpdateArrayCmd)
	DataCmd.AddCommand(RemoveArrayElementCmd)
}
