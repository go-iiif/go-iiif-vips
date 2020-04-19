package main

import (
	"context"
	_ "github.com/go-iiif/go-iiif-vips"
	"github.com/go-iiif/go-iiif/v2/tools"
	"log"
)

func main() {

	tool, err := tools.NewTransformTool()

	if err != nil {
		log.Fatal(err)
	}

	err = tool.Run(context.Background())

	if err != nil {
		log.Fatal(err)
	}
}
