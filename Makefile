cli-tools: 	
	go build -mod vendor -o bin/iiif-server cmd/iiif-server/main.go
	go build -mod vendor -o bin/iiif-tile-seed cmd/iiif-tile-seed/main.go
	go build -mod vendor -o bin/iiif-transform cmd/iiif-transform/main.go
	go build -mod vendor -o bin/iiif-process cmd/iiif-process/main.go

docker-build:
	docker build -f Dockerfile -t go-iiif-vips .

bump-version:
	perl -i -p -e 's/github.com\/go-iiif\/go-iiif\/$(PREVIOUS)/github.com\/go-iiif\/go-iiif\/$(NEW)/g' go.mod
	perl -i -p -e 's/github.com\/go-iiif\/go-iiif\/$(PREVIOUS)/github.com\/go-iiif\/go-iiif\/$(NEW)/g' README.md
	find . -name '*.go' | xargs perl -i -p -e 's/github.com\/go-iiif\/go-iiif\/$(PREVIOUS)/github.com\/go-iiif\/go-iiif\/$(NEW)/g'
