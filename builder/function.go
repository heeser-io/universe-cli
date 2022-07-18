package builder

import (
	"fmt"
	"log"
	"path"

	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
)

type UpdateAndUploadFunction struct {
	FunctionID  string
	Filepath    string
	Checksum    string
	Tags        []string
	Environment map[string]string
}

func ReleaseFunction(functionID string) error {
	releaseParams := v1.ReleaseFunctionParams{
		FunctionID: functionID,
	}

	if err := client.Client.Function.Release(&releaseParams); err != nil {
		return err
	}

	return nil
}

func UpdateFunction(params *v1.UpdateFunctionParams) (*v1.Function, error) {
	fileParams := v1.CreateFileParams{
		Filename:       path.Base(params.Path),
		Name:           fmt.Sprintf("function %s", params.FunctionID),
		Tags:           []string{"functions"},
		Checksum:       *params.Checksum,
		IsFunctionFile: true,
	}
	fileObj, err := client.Client.File.Create(&fileParams)
	if err != nil {
		log.Fatal(err)
	}

	if err := fileObj.Upload(params.Path); err != nil {
		log.Fatal(err)
	}
	fileObj.SignedUploadUrl = ""

	params.FileID = &fileObj.ID

	functionObj, err := client.Client.Function.Update(params)
	if err != nil {
		return nil, err
	}

	return functionObj, nil
}
func CreateFunction(params *v1.Function) (*v1.Function, error) {
	fileParams := v1.CreateFileParams{
		Filename:       path.Base(params.Path),
		Name:           fmt.Sprintf("function %s", params.Name),
		ProjectID:      params.ProjectID,
		Description:    fmt.Sprintf("file for the function %s", params.Name),
		Tags:           []string{"functions"},
		IsFunctionFile: true,
		Checksum:       params.Checksum,
	}
	fileObj, err := client.Client.File.Create(&fileParams)
	if err != nil {
		log.Fatal(err)
	}

	if err := fileObj.Upload(params.Path); err != nil {
		log.Fatal(err)
	}
	fileObj.SignedUploadUrl = ""

	createFunctionParams := v1.CreateFunctionParams{
		Name:        params.Name,
		ProjectID:   params.ProjectID,
		Checksum:    params.Checksum,
		Handler:     params.Handler,
		Permissions: params.Permissions,
		FileID:      fileObj.ID,
		Language:    params.Language,
		Tags:        params.Tags,
		Environment: params.Environment,
		BaseImage:   params.BaseImage,
	}

	functionObj, err := client.Client.Function.Create(&createFunctionParams)
	if err != nil {
		return nil, err
	}

	return functionObj, nil
}
