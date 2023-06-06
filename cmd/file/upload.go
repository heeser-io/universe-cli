package file

import (
	"fmt"
	"os"
	p "path"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	Path              string
	Level             string
	RandomizeFilename bool
	ProcessImages     bool
	UseExif           bool
	Reference         string
	AccountID         string
	Prefix            string
	Suffix            string
	Files             *[]string
	UploadCmd         = &cobra.Command{
		Use:   "upload",
		Short: "uploads a file",
		Run: func(cmd *cobra.Command, args []string) {
			files := map[string]*os.File{}

			for _, filePath := range *Files {
				file, err := os.Open(filePath)
				if err != nil {
					panic(err)
				}
				files[p.Base(filePath)] = file
			}

			uploadedFiles, err := client.Client.File.RawUpload(&v1.UploadAndCreateParams{
				Files:             files,
				Path:              Path,
				Level:             Level,
				RandomizeFilename: RandomizeFilename,
				ProcessImages:     ProcessImages,
				UseExif:           UseExif,
				Reference:         Reference,
				AccountID:         AccountID,
				Prefix:            Prefix,
				Suffix:            Suffix,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(uploadedFiles)))
		},
	}
)

func init() {
	Tags = UploadCmd.Flags().StringSlice("tags", nil, "tags of the file")
	Files = UploadCmd.Flags().StringSlice("files", nil, "paths to files")
	UploadCmd.Flags().BoolVar(&RandomizeFilename, "randomize-filename", false, "use to randomize all filenames")
	UploadCmd.Flags().BoolVar(&ProcessImages, "process-images", false, "use to process images on server to generate thumbnail and web size images")
	UploadCmd.Flags().BoolVar(&UseExif, "use-exif", false, "use to extract exif information of uploaded images (if available, for example for orientation)")
	UploadCmd.Flags().StringVar(&Path, "path", "", "use a custom path to store files in")
	UploadCmd.Flags().StringVar(&Level, "level", "private", "use private or public for file access")
	UploadCmd.Flags().StringVar(&Prefix, "prefix", "", "can be used to set a prefix to all filenames")
	UploadCmd.Flags().StringVar(&Suffix, "suffix", "", "cam ne ised to set a suffix to all filenames")
	UploadCmd.Flags().StringVar(&Reference, "reference", "", "used for filtering your files")
}
