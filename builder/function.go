package builder

import (
	"fmt"
	"log"
	"os"
	"path"

	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/joho/godotenv"
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

var (
	API_KEY string
)

func init() {
	godotenv.Load()
	API_KEY = os.Getenv("API_KEY")
}

func ReleaseFunction(functionID string) error {
	c := v1.WithAPIKey(API_KEY)

	releaseParams := v1.ReleaseFunctionParams{
		FunctionID: functionID,
	}

	if err := c.Function.Release(&releaseParams); err != nil {
		return err
	}

	return nil
}

func UpdateFunction(params *UpdateAndUploadFunction) (*v1.Function, error) {
	c := v1.WithAPIKey(API_KEY)

	fileParams := v1.CreateFileParams{
		Filename:       path.Base(params.Filepath),
		Name:           fmt.Sprintf("function %s", params.FunctionID),
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

	updateFunctionFileParams := v1.UpdateFunctionFileParams{
		FunctionID: params.FunctionID,
		FileID:     fileObj.ID,
		Checksum:   params.Checksum,
	}

	functionObj, err := c.Function.UpdateFile(&updateFunctionFileParams)
	if err != nil {
		return nil, err
	}

	return functionObj, nil
}
func CreateFunction(params *CreateAndUploadFunction) (*v1.Function, error) {
	c := v1.WithAPIKey(API_KEY)

	fileParams := v1.CreateFileParams{
		Filename:       path.Base(params.Filepath),
		Name:           fmt.Sprintf("function %s", params.Name),
		ProjectID:      params.ProjectID,
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

	createFunctionParams := v1.CreateFunctionParams{
		Name:      params.Name,
		ProjectID: params.ProjectID,
		Checksum:  params.Checksum,
		Handler:   params.Handler,
		FileID:    fileObj.ID,
		Language:  params.Language,
	}

	functionClient := v1.WithAPIKey(API_KEY)
	functionObj, err := functionClient.Function.Create(&createFunctionParams)
	if err != nil {
		return nil, err
	}

	return functionObj, nil
}
