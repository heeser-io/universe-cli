package stack

import (
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/builder"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	Force    bool
	Function string
	EnvFile  string
	PushCmd  = &cobra.Command{
		Use:   "push",
		Short: "pushes the current stack",
		Run: func(cmd *cobra.Command, args []string) {
			if EnvFile != "" {
				cwd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				if err := godotenv.Load(path.Join(cwd, EnvFile)); err != nil {
					panic(err)
				}
			}

			b, err := builder.New("", false)
			if err != nil {
				panic(err)
			}
			v := b.HasChange()

			stack, err := cmd.Flags().GetString("stack")
			if err != nil {
				panic(err)
			}

			if v || Force {
				if stack == "" {
					if err := b.BuildStack(Function); err != nil {
						color.Red("err: %v", err)
						os.Exit(1)
					}
					b.Verify()
				}

				subBuilders := b.GetSubBuilder()
				for _, subBuilder := range subBuilders {
					name := subBuilder.GetName()
					if name == stack || stack == "" {
						if err := subBuilder.BuildStack(Function); err != nil {
							fmt.Println(err)
						}
						subBuilder.Verify()
					}
				}
			}
		},
	}
)

func init() {
	PushCmd.Flags().BoolVarP(&Force, "force", "f", false, "force push")
	PushCmd.Flags().StringVar(&Function, "function", "", "only push this function in whole stack")
	PushCmd.Flags().StringVar(&EnvFile, "env-file", "", "relative path to env file for variable population")
}
