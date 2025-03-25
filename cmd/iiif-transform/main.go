package main

import (
	"context"
	"log"

	_ "github.com/aaronland/gocloud-blob/s3"
	_ "github.com/go-iiif/go-iiif-vips/v7"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	_ "gocloud.dev/blob/s3blob"

	"github.com/go-iiif/go-iiif/v6/app/transform"
)

func main() {

	ctx := context.Background()
	err = transform.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
