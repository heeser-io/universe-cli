package stack

import (
	"github.com/spf13/cobra"
)

var (
	Languages = []string{
		"golang:1.18.2",
	}
	Language    string
	SkeletonCmd = &cobra.Command{
		Use:   "skeleton",
		Short: "creates a skeleton project for a given language",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	SkeletonCmd.Flags().StringVarP(&Language, "language", "l", "", "programming language you want to use")
}
