package vips

// https://github.com/h2non/bimg
// https://github.com/jcupitt/libvips

import (
	"bytes"
	"errors"
	"fmt"
	iiifconfig "github.com/go-iiif/go-iiif/v3/config"
	iiifimage "github.com/go-iiif/go-iiif/v3/image"
	iiifsource "github.com/go-iiif/go-iiif/v3/source"
	"gopkg.in/h2non/bimg.v1"
	"image"
	"image/gif"
	_ "log"
)

type VIPSImage struct {
	iiifimage.Image
	config    *iiifconfig.Config
	source    iiifsource.Source
	source_id string
	id        string
	bimg      *bimg.Image
	isgif     bool
}

type VIPSDimensions struct {
	iiifimage.Dimensions
	imagesize bimg.ImageSize
}

func (d *VIPSDimensions) Height() int {
	return d.imagesize.Height
}

func (d *VIPSDimensions) Width() int {
	return d.imagesize.Width
}

/*

See notes in NewVIPSImageFromConfigWithSource - basically getting an image's
dimensions after the we've done the GIF conversion (just see the notes...)
will make bimg/libvips sad so account for that in Dimensions() and create a
pure Go implementation of the Dimensions interface (20160922/thisisaaronland)

*/

type GolangImageDimensions struct {
	iiifimage.Dimensions
	image image.Image
}

func (dims *GolangImageDimensions) Width() int {
	bounds := dims.image.Bounds()
	return bounds.Max.X
}

func (dims *GolangImageDimensions) Height() int {
	bounds := dims.image.Bounds()
	return bounds.Max.Y
}

func (im *VIPSImage) Update(body []byte) error {

	bimg := bimg.NewImage(body)
	im.bimg = bimg

	return nil
}

func (im *VIPSImage) Body() []byte {

	return im.bimg.Image()
}

func (im *VIPSImage) Format() string {

	return im.bimg.Type()
}

func (im *VIPSImage) ContentType() string {

	format := im.Format()

	if format == "jpg" || format == "jpeg" {
		return "image/jpeg"
	} else if format == "png" {
		return "image/png"
	} else if format == "webp" {
		return "image/webp"
	} else if format == "svg" {
		return "image/svg+xml"
	} else if format == "tif" || format == "tiff" {
		return "image/tiff"
	} else if format == "gif" {
		return "image/gif"
	} else {
		return ""
	}
}

func (im *VIPSImage) Identifier() string {
	return im.id
}

func (im *VIPSImage) Rename(id string) error {
	im.id = id
	return nil
}

func (im *VIPSImage) Dimensions() (iiifimage.Dimensions, error) {

	// see notes in NewVIPSImageFromConfigWithSource
	// ideally this never gets triggered but just in case...

	if im.isgif {

		buf := bytes.NewBuffer(im.Body())
		goimg, err := gif.Decode(buf)

		if err != nil {
			return nil, err
		}

		d := GolangImageDimensions{
			image: goimg,
		}

		return &d, nil
	}

	sz, err := im.bimg.Size()

	if err != nil {
		return nil, err
	}

	d := VIPSDimensions{
		imagesize: sz,
	}

	return &d, nil
}

// https://godoc.org/github.com/h2non/bimg#Options

func (im *VIPSImage) Transform(t *iiifimage.Transformation) error {

	// https://godoc.org/github.com/h2non/bimg#Options

	opts := bimg.Options{}
	opts.Quality = 100

	if t.Region != "full" {

		rgi, err := t.RegionInstructions(im)

		if err != nil {
			return err
		}

		if rgi.SmartCrop {

			opts.Gravity = bimg.GravitySmart
			opts.Crop = true
			opts.Width = rgi.Width
			opts.Height = rgi.Height

		} else {

			opts.AreaWidth = rgi.Width
			opts.AreaHeight = rgi.Height
			opts.Left = rgi.X
			opts.Top = rgi.Y

			if opts.Top == 0 && opts.Left == 0 {
				opts.Top = -1
			}

		}

		/*

			We need to do this or libvips will freak out and think it's trying to save
			an SVG file which it can't do (20160929/thisisaaronland)

		*/

		if im.ContentType() == "image/svg+xml" {
			opts.Type = bimg.PNG

		}

		/*
		   So here's a thing that we need to do because... computers?
		   (20160910/thisisaaronland)
		*/

		_, err = im.bimg.Process(opts)

		if err != nil {
			return err
		}

		// This is important and without it tiling gets completely
		// fubar-ed (20180620/thisisaaronland)
		// https://github.com/aaronland/go-iiif/issues/46

		opts = bimg.Options{}
		opts.Quality = 100

	} else {

		dims, err := im.Dimensions()

		if err != nil {
			return err
		}

		opts.Width = dims.Width()   // opts.AreaWidth,
		opts.Height = dims.Height() // opts.AreaHeight,
	}

	if t.Size != "max" && t.Size != "full" {

		dims, err := im.Dimensions()

		if err != nil {
			return err
		}

		width := dims.Width()
		height := dims.Height()

		si, err := t.SizeInstructionsWithDimensions(im, width, height)

		if err != nil {
			return err
		}

		opts.Width = si.Width
		opts.Height = si.Height
		opts.Force = si.Force
	}

	ri, err := t.RotationInstructions(im)

	if err != nil {
		return nil
	}

	// Okay, so there are a few things are happening here:

	// So apparently we need to explicitly swap height and width when we're auto-rotating
	// images based on EXIF orientation... which I guess makes but is kind of annoying but
	// so are a lot of things in life...

	// But it gets better... if you don't resize the image then all the metadata
	// gets preserved including the orientation flags which no longer make any
	// sense we've applied automagic orientation rotation hoohah so every application
	// that looks at the output (including this tool) will just keep rotating the
	// image by n-degrees every time. We could opts.StripMetaData all the images but that
	// seems a bit extreme and there is currently no way in bimg, libvips, some other pure-Go
	// package to strip or reset the Orientation EXIF tag. Instead we are hijacking the
	// IIIF spec to add make "-1" a valid rotation (see compliance/level2.go) which means
	// "do not autorotate based on EXIF orientation". As of this writing it does not work
	// in combination with other rotation commands (something like "-1,180" or "#180") but
	// it probably should... (20180607/thisisaaronland)

	// See also: https://github.com/h2non/bimg/issues/179

	if ri.NoAutoRotate {
		opts.NoAutoRotate = true
	} else {
		opts.Flip = ri.Flip
		opts.Rotate = bimg.Angle(ri.Angle % 360)
	}

	if !opts.NoAutoRotate {

		m, e := im.bimg.Metadata()

		if e == nil {

			// things that are on their side
			// https://magnushoff.com/jpeg-orientation.html

			if m.Orientation >= 5 && m.Orientation <= 8 {

				w := opts.Width
				h := opts.Height

				opts.Width = h
				opts.Height = w
			}
		}
	}

	if t.Quality == "color" || t.Quality == "default" {
		// do nothing.
	} else if t.Quality == "gray" {
		opts.Interpretation = bimg.InterpretationBW
	} else if t.Quality == "bitonal" {
		opts.Interpretation = bimg.InterpretationBW
	} else {
		// this should be trapped above
	}

	fi, err := t.FormatInstructions(im)

	if err != nil {
		return nil
	}

	if fi.Format == "jpg" {
		opts.Type = bimg.JPEG
	} else if fi.Format == "png" {
		opts.Type = bimg.PNG
	} else if fi.Format == "webp" {
		opts.Type = bimg.WEBP
	} else if fi.Format == "tif" {
		opts.Type = bimg.TIFF
	} else if fi.Format == "gif" {
		opts.Type = bimg.PNG // see this - we're just going to trick libvips until the very last minute...
	} else {
		msg := fmt.Sprintf("Unsupported image format '%s'", fi.Format)
		return errors.New(msg)
	}

	_, err = im.bimg.Process(opts)

	if err != nil {
		return err
	}

	err = iiifimage.ApplyCustomTransformations(t, im)

	if err != nil {
		return err
	}

	// see notes in NewVIPSImageFromConfigWithSource

	if fi.Format == "gif" && !im.isgif {

		goimg, err := iiifimage.IIIFImageToGolangImage(im)

		if err != nil {
			return err
		}

		im.isgif = true

		err = iiifimage.GolangImageToIIIFImage(goimg, im)

		if err != nil {
			return err
		}

	}

	return nil
}
