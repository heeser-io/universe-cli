package builder

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
)

type Filemapping struct {
	Name  string
	Path  string
	Files []*v1.File
	Tags  []string
}

func (fm *Filemapping) Upload(projectID string) ([]*v1.File, error) {
	files := []*v1.File{}

	fileInfo, err := os.Stat(fm.Path)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		if err := filepath.Walk(fm.Path, func(p string, info fs.FileInfo, err error) error {
			if !info.IsDir() {

				filename := path.Base(p)
				dir := path.Dir(p)

				f, err := os.Open(p)
				if err != nil {
					return err
				}

				cs, err := Checksum(p)
				if err != nil {
					return err
				}

				var eFile *v1.File
				for _, f := range fm.Files {
					if strings.Contains(f.Ref.Key, path.Join(p)) {
						eFile = f
					}
				}

				if eFile != nil && eFile.Checksum == cs {
					return nil
				}

				color.Yellow("uploading file %s...", p)
				fileRes, err := client.Client.File.RawUpload(&v1.UploadAndCreateParams{
					Path: dir,
					Files: map[string]*os.File{
						filename: f,
					},
					Level:     v1.LEVEL_PUBLIC,
					ProjectID: projectID,
					Tags:      fm.Tags,
				})
				if err != nil {
					return err
				}
				color.Green("successfully uploaded file %s", p)
				files = append(files, fileRes...)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	} else {
		f, err := os.Open(fm.Path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		fileRes, err := client.Client.File.RawUpload(&v1.UploadAndCreateParams{
			ProjectID: projectID,
			Files: map[string]*os.File{
				"file": f,
			},
		})
		if err != nil {
			return nil, err
		}
		files = append(files, fileRes...)
	}

	return files, nil
}
