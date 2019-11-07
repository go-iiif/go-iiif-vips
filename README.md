# go-iiif-vips

`go-iiif` driver for libvips.

## Important

This is work in progress. It should be considered to "work... until it doesn't".

## WIP

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

## See also

* https://github.com/go-iiif/go-iiif
* https://github.com/h2non/bimg/
* https://github.com/jcupitt/libvips
