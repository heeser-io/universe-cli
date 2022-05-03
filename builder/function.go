package builder

import (
	"fmt"
	"log"
	"os"
	"path"

	v1 "github.com/heeser-io/file-cdn/api-go/v1"
	functions_v1 "github.com/heeser-io/functions/api-go/v1"
	"github.com/joho/godotenv"
)

type CreateAndUploadFunction struct {
	Filepath string
	Name     string
	Language string
}

type UpdateAndUploadFunction struct {
	FunctionID string
	Filepath   string
}

var (
	API_KEY string
)

func init() {
	godotenv.Load()
	API_KEY = os.Getenv("API_KEY")
}
func ReleaseFunction(functionID string) error {
	c := functions_v1.WithAPIKey(API_KEY)

	releaseParams := functions_v1.ReleaseFunctionParams{
		FunctionID: functionID,
	}

	if err := c.Function.Release(&releaseParams); err != nil {
		return err
	}

	return nil
}

func UpdateFunction(params *UpdateAndUploadFunction) (*functions_v1.Function, error) {
	c := v1.WithAPIKey(API_KEY)

	fileParams := v1.CreateFileParams{
		Filename:       path.Base(params.Filepath),
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

	updateFunctionFileParams := functions_v1.UpdateFunctionFileParams{
		FunctionID: params.FunctionID,
		FileID:     fileObj.ID,
	}

	functionClient := functions_v1.WithAPIKey(API_KEY)
	functionObj, err := functionClient.Function.UpdateFile(&updateFunctionFileParams)
	if err != nil {
		return nil, err
	}

	return functionObj, nil
}
func CreateFunction(params *CreateAndUploadFunction) (*functions_v1.Function, error) {
	c := v1.WithAPIKey(API_KEY)

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

	functionClient := functions_v1.WithAPIKey(API_KEY)
	functionObj, err := functionClient.Function.Create(&createFunctionParams)
	if err != nil {
		return nil, err
	}

	return functionObj, nil
}
