package stack

import (
	"fmt"

	"github.com/heeser-io/universe-cli/builder"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

var (
	DownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "gives you the ability to download a project",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := builder.New("", true)
			if err != nil {
				panic(err)
			}

			projects, err := b.ListProjects()
			if err != nil {
				panic(err)
			}

			projectList := []string{}

			for _, project := range projects {
				projectList = append(projectList, fmt.Sprintf("%s (%s)", project.Name, project.ID))
			}

			prompt := promptui.Select{
				Label: "Select project to download",
				Items: projectList,
			}
			_, result, err := prompt.Run()
			if err != nil {
				panic(err)
			}

			selectedProject, ok := (funk.Find(projects, func(p v1.Project) bool {
				return fmt.Sprintf("%s (%s)", p.Name, p.ID) == result
			})).(v1.Project)

			if !ok {
				panic("unexpected error")
			}

			if err := b.Download(&selectedProject); err != nil {
				panic(err)
			}
		},
	}
)
