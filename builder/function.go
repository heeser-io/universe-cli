package builder

import (
	"fmt"
	"log"
	"path"

	v1 "github.com/heeser-io/universe/api/v1"
)

type CreateAndUploadFunction struct {
	ProjectID string
	Filepath  string
	Name      string
	Language  string
	Checksum  string
	Handler   string
}

type UpdateAndUploadFunction struct {
	FunctionID string
	Filepath   string
	Checksum   string
}

func ReleaseFunction(functionID string) error {
	releaseParams := v1.ReleaseFunctionParams{
		FunctionID: functionID,
	}

	if err := client.Function.Release(&releaseParams); err != nil {
		return err
	}

	return nil
}

func UpdateFunction(params *UpdateAndUploadFunction) (*v1.Function, error) {
	fileParams := v1.CreateFileParams{
		Filename:       path.Base(params.Filepath),
		Name:           fmt.Sprintf("function %s", params.FunctionID),
		Tags:           []string{"functions"},
		IsFunctionFile: true,
	}
	fileObj, err := client.File.Create(&fileParams)
	if err != nil {
		log.Fatal(err)
	}

	if err := fileObj.Upload(params.Filepath); err != nil {
		log.Fatal(err)
	}
	fileObj.SignedUploadUrl = ""

	updateFunctionFileParams := v1.UpdateFunctionFileParams{
		FunctionID: params.FunctionID,
		FileID:     fileObj.ID,
		Checksum:   params.Checksum,
	}

	functionObj, err := client.Function.UpdateFile(&updateFunctionFileParams)
	if err != nil {
		return nil, err
	}

	return functionObj, nil
}
func CreateFunction(params *CreateAndUploadFunction) (*v1.Function, error) {
	fileParams := v1.CreateFileParams{
		Filename:       path.Base(params.Filepath),
		Name:           fmt.Sprintf("function %s", params.Name),
		ProjectID:      params.ProjectID,
		Description:    fmt.Sprintf("file for the function %s", params.Name),
		Tags:           []string{"functions"},
		IsFunctionFile: true,
	}
	fileObj, err := client.File.Create(&fileParams)
	if err != nil {
		log.Fatal(err)
	}

	if err := fileObj.Upload(params.Filepath); err != nil {
		log.Fatal(err)
	}
	fileObj.SignedUploadUrl = ""

	createFunctionParams := v1.CreateFunctionParams{
		Name:      params.Name,
		ProjectID: params.ProjectID,
		Checksum:  params.Checksum,
		Handler:   params.Handler,
		FileID:    fileObj.ID,
		Language:  params.Language,
	}

	functionObj, err := client.Function.Create(&createFunctionParams)
	if err != nil {
		return nil, err
	}

	return functionObj, nil
}
