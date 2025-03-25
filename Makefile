CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli-tools: 	
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-server cmd/iiif-server/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-tile-seed cmd/iiif-tile-seed/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-transform cmd/iiif-transform/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-process cmd/iiif-process/main.go

docker-build:
	docker build -f Dockerfile -t go-iiif-vips .

bump-version:
	perl -i -p -e 's/github.com\/go-iiif\/go-iiif\/$(PREVIOUS)/github.com\/go-iiif\/go-iiif\/$(NEW)/g' go.mod
	perl -i -p -e 's/github.com\/go-iiif\/go-iiif\/$(PREVIOUS)/github.com\/go-iiif\/go-iiif\/$(NEW)/g' README.md
	find . -name '*.go' | xargs perl -i -p -e 's/github.com\/go-iiif\/go-iiif\/$(PREVIOUS)/github.com\/go-iiif\/go-iiif\/$(NEW)/g'
