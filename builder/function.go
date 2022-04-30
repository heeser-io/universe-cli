package builder

import (
	"fmt"
	"log"
	"os"
	"path"

	v1 "github.com/heeser-io/file-cdn/api-go/v1"
	functions_v1 "github.com/heeser-io/functions/api-go/v1"
)

type CreateAndUploaFunction struct {
	Filepath string
	Name     string
	Language string
}

func CreateFunction(params *CreateAndUploaFunction) (*functions_v1.Function, error) {
	apiKey := os.Getenv("API_KEY")

	c := v1.WithAPIKey(apiKey)

	fileParams := v1.CreateFileParams{
		Filename:       path.Base(params.Filepath),
		Description:    fmt.Sprintf("file for the function %s", params.Name),
		Tags:           []string{"functions"},
		IsFunctionFile: true,
	}
	fileObj, err := c.File.Create(&fileParams)
	if err != nil {
		log.Fatal(err)
	}

	if err := fileObj.Upload(params.Filepath); err != nil {
		log.Fatal(err)
	}
	fileObj.SignedUploadUrl = ""

	createFunctionParams := functions_v1.CreateFunctionParams{
		Name:     params.Name,
		FileID:   fileObj.ID,
		Language: params.Language,
	}

	functionClient := functions_v1.WithAPIKey(apiKey)
	functionObj, err := functionClient.Function.Create(&createFunctionParams)
	if err != nil {
		return nil, err
	}

	return functionObj, nil
}
