# go-iiif-vips

`go-iiif` driver for libvips.

## Important

You should start by reading the documentation in the [go-iiif](https://github.com/go-iiif/go-iiif/blob/master/README.md) package.

## Tools

As of version 2 of [go-iiif](https://github.com/go-iiif/go-iiif) all of the logic, including defining and parsing command line arguments, for any `go-iiif` tool that performs image processing has been moved in to the `tools` package. This change allows non-core image processing packages (like [go-iiif-vips](https://github.com/go-iiif/go-iiif-vips)) to more easily re-use functionality defined in the core `go-iiif` package. For example:

```
package main

import (
	"context"
	_ "github.com/go-iiif/go-iiif-vips"
	"github.com/go-iiif/go-iiif/tools"
)

func main() {
	tool, _ := tools.NewProcessTool()
	tool.Run(context.Background())
}
```

### iiif-process

```
$> ./bin/iiif-process -h
Usage of ./bin/iiif-process:
  -config string
    	Path to a valid go-iiif config file. DEPRECATED - please use -config_source and -config name.
  -config-name string
    	The name of your go-iiif config file. (default "config.json")
  -config-source string
    	A valid Go Cloud bucket URI where your go-iiif config file is located.
  -instructions string
    	Path to a valid go-iiif processing instructions file. DEPRECATED - please use -instructions-source and -instructions-name.
  -instructions-name string
    	The name of your go-iiif instructions file. (default "instructions.json")
  -instructions-source string
    	A valid Go Cloud bucket URI where your go-iiif instructions file is located.
  -mode string
    	Valid modes are: cli, lambda. (default "cli")
  -report
    	Store a process report (JSON) for each URI in the cache tree.
  -report-name string
    	The filename for process reports. Default is 'process.json' as in '${URI}/process.json'. (default "process.json")
```

Perform a series of IIIF image processing tasks, defined in a JSON-based "instructions" file, on one or more (IIIF) URIs. For example:

```
$> ./bin/iiif-process -config config.json -instructions instructions.json -uri source/IMG_0084.JPG | jq

{
  "source/IMG_0084.JPG": {
    "dimensions": {
      "b": [
        2048,
        1536
      ],
      "d": [
        320,
        320
      ],
      "o": [
        4032,
        3024
      ]
    },
    "palette": [
      {
        "name": "#b87531",
        "hex": "#b87531",
        "reference": "vibrant"
      },
      {
        "name": "#805830",
        "hex": "#805830",
        "reference": "vibrant"
      },
      {
        "name": "#7a7a82",
        "hex": "#7a7a82",
        "reference": "vibrant"
      },
      {
        "name": "#c7c3b3",
        "hex": "#c7c3b3",
        "reference": "vibrant"
      },
      {
        "name": "#5c493a",
        "hex": "#5c493a",
        "reference": "vibrant"
      }
    ],
    "uris": {
      "b": "source/IMG_0084.JPG/full/!2048,1536/0/color.jpg",
      "d": "source/IMG_0084.JPG/-1,-1,320,320/full/0/dither.jpg",
      "o": "source/IMG_0084.JPG/full/full/-1/color.jpg"
    }
  }
}
```

Images are read-from and stored-to whatever source or derivatives caches defined in your `config.json` file.

#### "instructions" files

An instruction file is a JSON-encoded dictionary. Keys are user-defined and values are dictionary of IIIF one or more transformation instructions. For example:

```
{
    "o": {"size": "full", "format": "", "rotation": "-1" },
    "b": {"size": "!2048,1536", "format": "jpg" },
    "d": {"size": "full", "quality": "dither", "region": "-1,-1,320,320", "format": "jpg" }	
}

```

The complete list of possible instructions is:

```
type IIIFInstructions struct {
	Region   string `json:"region"`
	Size     string `json:"size"`
	Rotation string `json:"rotation"`
	Quality  string `json:"quality"`
	Format   string `json:"format"`
}
```

As of this writing there is no explicit response type for image beyond `map[string]interface{}`. There probably could be but it's still early days.

### iiif-server

```
$> ./bin/iiif-server -h
Usage of ./bin/iiif-server:
  -config string
    	Path to a valid go-iiif config file. DEPRECATED - please use -config-url and -config name.
  -config-name string
    	The name of your go-iiif config file. (default "config.json")
  -config-source string
    	A valid Go Cloud bucket URI where your go-iiif config file is located.	
  -example
    	Add an /example endpoint to the server for testing and demonstration purposes
  -example-root string
    	An explicit path to a folder containing example assets (default "example")
  -host string
    	Bind the server to this host (default "localhost")
  -port int
    	Bind the server to this port (default 8080)
  -protocol string
    	The protocol for wof-staticd server to listen on. Valid protocols are: http, lambda. (default "http")
```

```
$> bin/iiif-server -config config.json
2016/09/01 15:45:07 Serving 127.0.0.1:8080 with pid 12075

curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/full/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/125,15,200,200/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/pct:41.6,7.5,40,70/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/full/full/270/default.png
```

`iiif-server` is a HTTP server that supports version 2.1 of the [IIIF Image API](http://iiif.io/api/image/2.1/).

#### Endpoints

Although the identifier parameter (`{ID}`) in the examples below suggests that is is only string characters up to and until a `/` character, it can in fact contain multiple `/` separated strings. For example, either of these two URLs is valid

```
http://localhost:8082/191733_5755a1309e4d66a7_k.jpg/info.json
http://localhost:8082/191/733/191733_5755a1309e4d66a7/info.json
```

Where the identified will be interpreted as `191733_5755a1309e4d66a7_k.jpg` and `191/733/191733_5755a1309e4d66a7` respectively. Identifiers containing one or more `../` strings will be made to feel bad about themselves.

##### GET /{ID}/info.json

```
$> curl -s http://localhost:8082/184512_5f7f47e5b3c66207_x.jpg/info.json | python -mjson.tool
{
    "@context": "http://iiif.io/api/image/2/context.json",
    "@id": "http://localhost:8082/184512_5f7f47e5b3c66207_x.jpg",
    "@type": "iiif:Image",
    "height": 4096,
    "profile": [
        "http://iiif.io/api/image/2/level2.json",
        {
            "formats": [
                "tif",
                "webp",
                "jpg",
                "png"
            ],
            "qualities": [
                "default",
		"dither",
                "color"
            ],
            "supports": [
                "full",
                "regionByPx",
                "regionByPct",
                "sizeByWh",
                "full",
                "max",
                "sizeByW",
                "sizeByH",
                "sizeByPct",
                "sizeByConfinedWh",
                "none",
                "rotationBy90s",
                "mirroring",
                "baseUriRedirect",
                "cors",
                "jsonldMediaType"
            ]
        }
    ],
    "protocol": "http://iiif.io/api/image",
    "width": 3897
}
```

Return the [profile description](http://iiif.io/api/image/2.1/#profile-description) for an identifier.

##### GET /{ID}/{REGION}/{SIZE}/{ROTATION}/{QUALITY}.{FORMAT}

```
$> curl -s http://localhost:8082/184512_5f7f47e5b3c66207_x.jpg/pct:41,7,40,70/,250/0/default.jpg
```

Return an image derived from an identifier and one or more [IIIF parameters](http://iiif.io/api/image/2.1/#image-request-parameters). For example:

![spanking cat, cropped](misc/go-iiif-crop.jpg)

##### GET /debug/vars

```
$> curl -s 127.0.0.1:8080/debug/vars | python -mjson.tool | grep Cache
    "CacheHit": 4,
    "CacheMiss": 16,
    "CacheSet": 16,

$> curl -s 127.0.0.1:8080/debug/vars | python -mjson.tool | grep Transforms
    "TransformsAvgTimeMS": 1833.875,
    "TransformsCount": 16,
```

This exposes all the usual Go [expvar](https://golang.org/pkg/expvar/) debugging output along with the following additional properies:

* CacheHit - _the total number of (derivative) images successfully returned from cache_
* CacheMiss - _the total number of (derivative) images not found in the cache_
* CacheSet - _the total number of (derivative) images added to the cache_
* TransformsAvgTimeMS - _the average amount of time in milliseconds to transforms a source image in to a derivative_
* TransformsCount - _the total number of source images transformed in to a derivative_

_Note: This endpoint is only available from the machine the server is running on._

### iiif-tile-seed

```
$> ./bin/iiif-tile-seed -h
Usage of ./bin/iiif-tile-seed:
  -config string
    	Path to a valid go-iiif config file. DEPRECATED - please use -config-source and -config name.
  -config-name string
    	The name of your go-iiif config file. (default "config.json")
  -config-source string
    	A valid Go Cloud bucket URI where your go-iiif config file is located.
  -csv-source string
    	 (default "A valid Go Cloud bucket URI where your CSV tileseed files are located.")
  -endpoint string
    	The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image (default "http://localhost:8080")
  -format string
    	A valid IIIF format parameter (default "jpg")
  -logfile string
    	Write logging information to this file
  -loglevel string
    	The amount of logging information to include, valid options are: debug, info, status, warning, error, fatal (default "info")
  -mode string
    	Valid modes are: cli, csv, lambda. (default "cli")
  -noextension
    	Remove any extension from destination folder name.
  -processes int
    	The number of concurrent processes to use when tiling images (default 4)
  -quality string
    	A valid IIIF quality parameter - if "default" then the code will try to determine which format you've set as the default (default "default")
  -refresh
    	Refresh a tile even if already exists (default false)
  -scale-factors string
    	A comma-separated list of scale factors to seed tiles with (default "4")
  -verbose
    	Write logging to STDOUT in addition to any other log targets that may have been defined
```

Generate (seed) all the tiled derivatives for a source image for use with the [Leaflet-IIIF](https://github.com/mejackreed/Leaflet-IIIF) plugin.

#### iiif-tile-seed and identifiers

Identifiers for source images can be passed to `iiif-tiles-seed` in of two way:

1. A space-separated list of identifiers
2. A space-separated list of _comma-separated_ identifiers indicating the identifier for the source image followed by the identifier for the newly generated tiles

For example:

```
$> ./bin/iiif-tile-seed -options 191733_5755a1309e4d66a7_k.jpg
```

Or:

```
$> ./bin/iiif-tile-seed -options 191733_5755a1309e4d66a7_k.jpg,191/733/191733_5755a1309e4d66a7
```

In many cases the first option will suffice but sometimes you might need to create new identifiers or structure existing identifiers according to their output, for example avoiding the need to store lots of file in a single directory. It's up to you.

You can also run `iiif-tile-seed` pass a list of identifiers as a CSV file. To do so include the `-mode csv` argument, like this:

```
$> ./bin/iiif-tile-seed -options -mode csv CSVFILE
```

Your CSV file must contain a header specifying a `source_id` and `alternate_id` column, like this:

```
source_id,alternate_id
191733_5755a1309e4d66a7_k.jpg,191733_5755a1309e4d66a7
```

While all columns are required if `alternate_id` is empty the code will simply default to using `source_id` for all operations.

_Important: The use of alternate IDs is not fully supported by `iiif-server` yet. Which is to say to the logic for how to convert a source identifier to an alternate identifier is still outside the scope of `go-iiif` so unless you have pre-rendered all of your tiles or other derivatives (in which case the check for cached derivatives at the top of the imgae handler will be triggered) then the server won't know where to write new alternate files._

## Config files

You should start by reading the [documentation for configuation files](https://github.com/go-iiif/go-iiif/blob/master/README.md#config-files) in the `go-iiif` package. What follows are configuration options specific to the `go-iiif-vips` package.

### graphics

```
	"graphics": {
		"source": { "name": "vips" }
	}
```

According to the [bimg docs](https://github.com/h2non/bimg/) (which is the Go library wrapping `libvips`) the following formats can be read:

```
It can read JPEG, PNG, WEBP natively, and optionally TIFF, PDF, GIF and SVG formats if libvips@8.6+ is compiled with proper library bindings.
```

If you've installed `libvips` using [the handy setup script](setup/setup-libvips-ubuntu.sh) then all the formats listed above, save PDF, [should be supported](https://github.com/jcupitt/libvips#optional-dependencies).

_Important: That's actually not true if you're reading this. It was true but then I tried running `iiif-tile-seed` on a large set of images and started triggering [this error](https://github.com/h2non/bimg/issues/111) even though it's supposed to be fixed. If you're reading this it means at least one of three things: the bug still exists; I pulled source from `gopkg.in` rather than `github.com` despite the author's notes in the issue; changes haven't propogated to `gopkg.in` yet. Which is to say that the current version of `bimg` is pegged to the [v1.0.1](https://github.com/h2non/bimg/releases/tag/v1.0.1) release which doesn't know think it knows about the PDF, GIF or SVG formats yet. It's being worked on..._

The `VIPS` graphics source has the following optional properties:

* **tmpdir** Specify an alternate path where libvips [should write temporary files](http://www.vips.ecs.soton.ac.uk/supported/7.42/doc/html/libvips/VipsImage.html#vips-image-new-temp-file) while processing images. This may be necessary if you are a processing many large files simultaneously and your default "temporary" directory is very small.

## Docker

Yes. There is a [Dockerfile](Dockerfile) included with this distribution. It will build a container with the following tools:

* The `iiif-server` tool.
* The `iiif-process` command-line tool.
* The `iiif-tile-seed` command-line tool.

To build the container run:

```
$> docker build -f Dockerfile -t go-iiif-vips .
```

To start the `iiif-server` tool run:

```
$> docker run -it -p 6161:8080 \
   -v /usr/local/go-iiif/docker/etc:/etc/iiif-server \
   -v /usr/local/go-iiif/docker/images:/usr/local/iiif-server \
   go-iiif-vips \
   /bin/iiif-server -host 0.0.0.0 -config /etc/iiif-server/config.json
   
2018/06/20 23:03:10 Listening for requests at 0.0.0.0:8080
```

See the way we are mapping `/etc/iiif-server` and `/usr/local/iiif-server` to local directories? By default the `iiif-server` Dockerfile does not bundle config files or images. Maybe some day, but that day is not today.

Then, in another terminal:

```
$> curl localhost:6161/test.jpg/info.json
{"@context":"http://iiif.io/api/image/2/context.json","@id":"http://localhost:6161/test.jpg","@type":"iiif:Image","protocol":"http://iiif.io/api/image","width":3897,"height":4096,"profile":["http://iiif.io/api/image/2/level2.json",{"formats":["gif","webp","jpg","png","tif"],"qualities":["default","color","dither"],"supports":["full","regionByPx","regionByPct","regionSquare","sizeByDistortedWh","sizeByWh","full","max","sizeByW","sizeByH","sizeByPct","sizeByConfinedWh","none","rotationBy90s","mirroring","noAutoRotate","baseUriRedirect","cors","jsonldMediaType"]}],"service":[{"@context":"x-urn:service:go-iiif#palette","profile":"x-urn:service:go-iiif#palette","label":"x-urn:service:go-iiif#palette","palette":[{"name":"#2f2013","hex":"#2f2013","reference":"vibrant"},{"name":"#9e8e65","hex":"#9e8e65","reference":"vibrant"},{"name":"#c6bca6","hex":"#c6bca6","reference":"vibrant"},{"name":"#5f4d32","hex":"#5f4d32","reference":"vibrant"}]}]}
```

Let's say you're using S3 as an image source and reading (S3) credentials from environment variables (something like `{"source": { "name": "S3", "path": "{BUCKET}", "region": "us-east-1", "credentials": "env:" }`) then you would start up `iiif-server` like this:

```
$> docker run -it -p 6161:8080 \
       -v /usr/local/go-iiif/docker/etc:/etc/iiif-server -v /usr/local/go-iiif/docker/images:/usr/local/iiif-server \
       -e AWS_ACCESS_KEY_ID={AWS_KEY} -e AWS_SECRET_ACCESS_KEY={AWS_SECRET} \
       go-iiif-vips-server \
       /bin/iiif-server -host 0.0.0.0 -config /etc/iiif-server/config.json       
```

The process an image using the `iiif-process` Docker container you would run something like:

```
$> docker run \
   -v /usr/local/go-iiif/docker/etc:/etc/go-iiif \
   go-iiif-vips \
   /bin/iiif-process -config=/etc/go-iiif/config.json -instructions=/etc/go-iiif/instructions.json \
   -uri=test.jpg
```

Again, see the way we're mapping `/etc/go-iiif` to a local folder, like we do in the `iiif-server` Docker example? The same rules apply here.

## See also

* https://github.com/go-iiif/go-iiif
* https://github.com/h2non/bimg/
* https://github.com/jcupitt/libvips
