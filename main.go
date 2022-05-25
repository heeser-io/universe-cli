package main

import (
	"github.com/heeser-io/universe-cli/cmd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}
