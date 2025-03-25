package vips

import (
	"context"

	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifimage "github.com/go-iiif/go-iiif/v6/image"
	iiifsource "github.com/go-iiif/go-iiif/v6/source"
	"github.com/h2non/bimg"
)

func init() {

	dr, err := NewVIPSDriver()

	if err != nil {
		panic(err)
	}

	iiifdriver.RegisterDriver("vips", dr)
}

type VIPSDriver struct {
	iiifdriver.Driver
}

func NewVIPSDriver(ctx context.Context, uri string) (iiifdriver.Driver, error) {
	dr := &VIPSDriver{}
	return dr, nil
}

func (dr *VIPSDriver) NewImageFromConfigWithSource(ctx context.Context, config *iiifconfig.Config, src iiifsource.Source, id string) (iiifimage.Image, error) {

	body, err := src.Read(id)

	if err != nil {
		return nil, err
	}

	bimg := bimg.NewImage(body)

	im := VIPSImage{
		config:    config,
		source:    src,
		source_id: id,
		id:        id,
		bimg:      bimg,
		isgif:     false,
	}

	/*

		Hey look - see the 'isgif' flag? We're going to hijack the fact that
		bimg doesn't handle GIF files and if someone requests them then we
		will do the conversion after the final call to im.bimg.Process and
		after we do handle any custom features. We are relying on the fact
		that both bimg.NewImage and bimg.Image() expect and return raw bytes
		and we are ignoring whatever bimg thinks in the Format() function.
		So basically you should not try to any processing in bimg/libvips
		after the -> GIF transformation. (20160922/thisisaaronland)

		See also: https://github.com/h2non/bimg/issues/41
	*/

	return &im, nil
}

func (dr *VIPSDriver) NewImageFromConfigWithCache(ctx context.Context, config *iiifconfig.Config, cache iiifcache.Cache, id string) (iiifimage.Image, error) {

	var image iiifimage.Image

	body, err := cache.Get(id)

	if err == nil {

		source, err := iiifsource.NewMemorySource(body)

		if err != nil {
			return nil, err
		}

		image, err = dr.NewImageFromConfigWithSource(ctx, config, source, id)

		if err != nil {
			return nil, err
		}

	} else {

		image, err = dr.NewImageFromConfig(ctx, config, id)

		if err != nil {
			return nil, err
		}

		go func() {
			cache.Set(id, image.Body())
		}()
	}

	return image, nil
}

func (dr *VIPSDriver) NewImageFromConfig(ctx context.Context, config *iiifconfig.Config, id string) (iiifimage.Image, error) {

	source, err := iiifsource.NewSource(ctx, config.Images.Source.URI)

	if err != nil {
		return nil, err
	}

	return dr.NewImageFromConfigWithSource(config, source, id)
}
