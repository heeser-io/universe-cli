package main

import (
	"github.com/heeser-io/stack-builder/builder"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	builder.BuildStack("./universe.yml")
}
